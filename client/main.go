package main

import (
	"context"
	"fmt"
	"log"

	"github.com/btwiuse/h3/utils"
	"github.com/marten-seemann/webtransport-go"
)

func handleConn(conn *webtransport.Session) {
	log.Println("new conn")
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
