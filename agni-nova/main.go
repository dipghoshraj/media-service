package main

import (
	"net"
	"os"

	"github.com/odio4u/agni-tunnels/agni-nova/nova"
)

func main() {
	port := nova.YamlConfig.Nova.Port
	if port == "" {
		port = "3001"
	}
	nova.Logger.Info("starting nova proxy", "listen_addr", ":"+port)

	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		nova.Logger.Error("failed to start listener", "listen_addr", ":"+port, "error", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			nova.Logger.Error("failed to accept connection", "error", err)
			continue
		}
		go nova.HandleStream(conn)
	}
}
