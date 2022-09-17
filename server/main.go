package server

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/btwiuse/h3/utils"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/marten-seemann/webtransport-go"
)

func makeServer(host, port, cert, key string) *Server {
	server := webtransport.Server{
		H3: http3.Server{
			Addr: port,
		},
	}
	return &Server{
		Server: server,
		Host:   host,
		Port:   port,
		Cert:   cert,
		Key:    key,
	}
}

type Server struct {
	webtransport.Server

	Host string
	Port string
	Cert string
	Key  string
}

func (s *Server) ListenAndServeTLS() error {
	log.Printf("listening on https://%s%s (UDP)", s.Host, s.Port)
	err := s.Server.ListenAndServeTLS(s.Cert, s.Key)
	log.Fatalln(err)
	return err
}

func (s *Server) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc("/echo", handler)
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

func echoConn(conn *webtransport.Session) {
	log.Println("new conn", conn.RemoteAddr())
	stream, err := conn.AcceptStream(context.TODO())
	if err != nil {
		log.Println("error accepting stream")
		return
	}
	io.Copy(stream, stream)
}

func Run([]string) error {
	s := makeServer(
		utils.EnvHOST("localhost"),
		utils.EnvPORT(":443"),
		utils.EnvCERT("localhost.pem"),
		utils.EnvKEY("localhost-key.pem"),
	)
	s.HandleFunc("/echo", s.handleEcho)
	return s.ListenAndServeTLS()
}
