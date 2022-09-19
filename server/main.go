package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/btwiuse/h3/utils"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
	"k0s.io/pkg/reverseproxy"
)

func makeServer(host, port, altsvc, cert, key string) *Server {
	server := webtransport.Server{
		H3: http3.Server{
			Addr: port,
		},
		CheckOrigin: func(*http.Request) bool {
			return true
		},
	}
	return &Server{
		Server: server,
		Host:   host,
		Port:   port,
		AltSvc: altsvc,
		Cert:   cert,
		Key:    key,
	}
}

type Server struct {
	webtransport.Server

	Host   string
	Port   string
	AltSvc string
	Cert   string
	Key    string
}

func (s *Server) ListenAndServeTLS() error {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", s.handleHTTP1)
		log.Printf("listening on http://%s%s (TCP)", s.Host, s.Port)
		http.ListenAndServe(s.Port, mux)
	}()
	log.Printf("listening on https://%s%s (UDP)", s.Host, s.Port)
	err := s.Server.ListenAndServeTLS(s.Cert, s.Key)
	log.Fatalln(err)
	return err
}

func (s *Server) Handle(path string, handler http.Handler) {
	http.Handle(path, handler)
}

func (s *Server) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(path, handler)
}

func (s *Server) handleHTTP1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", s.AltSvc)
	http.Error(w, fmt.Sprintf("Alt-Svc: %s", s.AltSvc), 200)
}

func (s *Server) handleEcho(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrade(w, r)
	if err != nil {
		log.Printf("upgrading failed: %s", err)
		w.WriteHeader(500)
		return
	}
	go echoConn(conn)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	uiURL := "https://http3.vercel.app/"
	reverseproxy.Handler(uiURL).ServeHTTP(w, r)
}

func echoConn(conn *webtransport.Session) {
	log.Println(conn.RemoteAddr(), "new session")
	ctx := context.Background()
	for {
		stream, err := conn.AcceptStream(ctx)
		if err != nil {
			log.Println(conn.RemoteAddr(), "session closed")
			break
		}
		log.Println(conn.RemoteAddr(), "new stream")
		go io.Copy(stream, stream)
	}
}

func ApplyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("[H3]", r.RemoteAddr, "->", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func Run([]string) error {
	var (
		host   = utils.EnvHost("localhost")
		port   = utils.EnvPort(":443")
		altsvc = utils.EnvAltSvc(fmt.Sprintf(`h3="%s"`, port))
		cert   = utils.EnvCert("localhost.pem")
		key    = utils.EnvKey("localhost-key.pem")
		s      = makeServer(host, port, altsvc, cert, key)
	)
	s.Handle("/", http.HandlerFunc(s.handleRoot))
	s.Handle("/echo", ApplyMiddleware(http.HandlerFunc(s.handleEcho)))
	return s.ListenAndServeTLS()
}
