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

func NewServer(host, port, altsvc, cert, key string) *Server {
	s := &Server{
		Host:   host,
		Port:   port,
		AltSvc: altsvc,
		Cert:   cert,
		Key:    key,
	}
	s.server = s.webtransportServer()
	return s
}

type Server struct {
	server *webtransport.Server

	Host   string
	Port   string
	AltSvc string
	Cert   string
	Key    string
}

func (s *Server) ListenAndServe() error {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", s.handleHTTP1)
		log.Printf("listening on http://%s%s (TCP)", s.Host, s.Port)
		http.ListenAndServe(s.Port, mux)
	}()
	log.Printf("listening on https://%s%s (UDP)", s.Host, s.Port)
	return s.server.ListenAndServeTLS(s.Cert, s.Key)
}

func (s *Server) handleHTTP1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Alt-Svc", s.AltSvc)
	http.Error(w, fmt.Sprintf("Alt-Svc: %s", s.AltSvc), 200)
}

func (s *Server) handleEcho(w http.ResponseWriter, r *http.Request) {
	conn, err := s.server.Upgrade(w, r)
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

func (s *Server) webtransportServer() *webtransport.Server {
	return &webtransport.Server{
		H3: http3.Server{
			Addr:            s.Port,
			Handler:         s.handler(),
			EnableDatagrams: true,
		},
		CheckOrigin: func(*http.Request) bool { return true },
	}
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(s.handleRoot))
	mux.Handle("/echo", applyMiddleware(http.HandlerFunc(s.handleEcho)))
	return mux
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

func applyMiddleware(next http.Handler) http.Handler {
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
		s      = NewServer(host, port, altsvc, cert, key)
	)
	return s.ListenAndServe()
}
