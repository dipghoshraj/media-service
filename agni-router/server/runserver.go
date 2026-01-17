package server

import (
	"log"
	"net"

	"github.com/Purple-House/agni-tunnels/agni-router/pkg/session"
)

func RouterServer() {
	ln, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go session.HandleStream(conn)
	}
}
