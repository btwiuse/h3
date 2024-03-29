package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/btwiuse/h3/utils"
	"github.com/quic-go/webtransport-go"
)

func Run([]string) error {
	u := fmt.Sprintf(
		"https://%s%s/echo",
		utils.EnvHost("localhost"),
		utils.EnvPort(":443"),
	)
	ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	var d webtransport.Dialer
	log.Printf("dialing %s (UDP)", u)
	resp, conn, err := d.Dial(ctx, u, nil)
	if err != nil {
		log.Fatalln(err)
	}
	_ = resp
	handleConn(conn)
	return nil
}

func handleConn(conn *webtransport.Session) {
	log.Println("new conn", conn.LocalAddr())
	ctx, _ := context.WithTimeout(context.TODO(), time.Second)
	stream, err := conn.OpenStreamSync(ctx)
	if err != nil {
		log.Println("error opening stream:", err)
		return
	}
	go io.Copy(os.Stdout, stream)
	io.Copy(stream, os.Stdin)
}
