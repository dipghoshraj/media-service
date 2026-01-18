package main

import (
	"log"
	"net"

	"github.com/odio4u/agni-tunnels/agni-nova/nova"
)

func main() {
	log.Println("This is the main entry point for the indraNet reverse proxy. Port :---", 3000)

	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go nova.HandleStream(conn)
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
	}
}
