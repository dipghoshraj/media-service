package server

import (
	"fmt"
	"net"

	"github.com/odio4u/agni-tunnels/agni-router/pkg/config"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/session"
)

func RouterServer() {
	port := fmt.Sprintf(":%s", config.YamlConfig.Router.ProxtPort)
	logger.Logger.Info("TCP proxy listener starting", "listen_addr", port)
	ln, err := net.Listen("tcp", port)
	if err != nil {
		logger.Logger.Error("failed to start TCP listener", "listen_addr", port, "error", err)
		return
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Logger.Error("failed to accept connection", "error", err)
			continue
		}
		go session.HandleStream(conn)
	}
}
