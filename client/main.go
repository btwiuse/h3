package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/btwiuse/h3/utils"
	"github.com/marten-seemann/webtransport-go"
)

func handleConn(conn *webtransport.Session) {
	log.Println("new conn", conn.LocalAddr())
	stream, err := conn.OpenStream()
	if err != nil {
		log.Println("error opening stream:", err)
	}
	go io.Copy(os.Stdout, stream)
	io.Copy(stream, os.Stdin)
}

func main() {
	var d webtransport.Dialer
	resp, conn, err := d.Dial(context.TODO(), fmt.Sprintf("https://localhost%s/echo", utils.EnvPORT(":443")), nil)
	if err != nil {
		log.Fatalln(err)
	}
	_ = resp
	handleConn(conn)
}
