package server

import (
	"fmt"
	"log"
	"net"

	"github.com/odio4u/agni-tunnels/agni-router/pkg/config"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/session"
)

func RouterServer() {
	port := fmt.Sprintf(":%s", config.YamlConfig.Router.ProxtPort)
	log.Println("Running server on ", port)
	ln, err := net.Listen("tcp", port)
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
