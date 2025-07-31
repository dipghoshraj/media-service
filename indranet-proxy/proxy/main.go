package main

import (
	"cosmo-proxy/internal"
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("This is the main entry point for the indraNet reverse proxy.")

	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go internal.HandleConnection(conn)
		log.Printf("Accepted connection from %s", conn.RemoteAddr())
	}
}
