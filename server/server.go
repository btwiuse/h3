package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/btwiuse/h3/utils"
	"github.com/webtransport/quic-go/http3"
	"github.com/webtransport/webtransport-go"
	"k0s.io/pkg/reverseproxy"
)

func NewServer(host, port, altsvc, ui, cert, key string) *Server {
	s := &Server{
		Host:   host,
		Port:   port,
		AltSvc: altsvc,
		UI:     ui,
		Cert:   cert,
		Key:    key,
	}
	s.server = s.webtransportServer()
	s.uiHandler = reverseproxy.Handler(s.UI)
	return s
}

type Server struct {
	server    *webtransport.Server
	uiHandler http.Handler

	Host   string
	Port   string
	AltSvc string
	UI     string
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
	s.uiHandler.ServeHTTP(w, r)
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
	s.uiHandler.ServeHTTP(w, r)
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
		ui     = utils.EnvUI("https://http3.vercel.app")
		cert   = utils.EnvCert("localhost.pem")
		key    = utils.EnvKey("localhost-key.pem")
		s      = NewServer(host, port, altsvc, ui, cert, key)
	)
	return s.ListenAndServe()
}
