package session

import (
	"net"
	"sync"

	"github.com/google/uuid"
	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
	"github.com/odio4u/mem-sdk/sni"
)

type TunnleContext struct {
	connection_id string
	stream        *tunnel.AgniTunnel_ConnectServer
	tcp           net.Conn
	closed        chan struct{}
	closeOnce     sync.Once
}

func HandleStream(conn net.Conn) {
	clientAddr := conn.RemoteAddr().String()
	logger.Logger.Info("connection accepted", "client_addr", clientAddr)

	serverName, wrappedConn, err := sni.PeekSNI(conn)
	if err != nil {
		logger.Logger.Error("failed to extract SNI", "client_addr", clientAddr, "error", err)
		conn.Close()
		return
	}
	logger.Logger.Info("SNI extracted", "sni", serverName, "client_addr", clientAddr)

	session, exists := Seeder.GetSession(serverName)
	if !exists {
		logger.Logger.Warn("no active session for SNI", "sni", serverName, "client_addr", clientAddr)
		conn.Close()
		return
	}

	tunnelContext := &TunnleContext{
		connection_id: uuid.New().String(),
		stream:        session.Stream,
		tcp:           wrappedConn,
		closed:        make(chan struct{}),
	}
	logger.Logger.Info("tunnel context created", "sni", serverName, "connection_id", tunnelContext.connection_id)

	sendOpen(tunnelContext)

	go PollGRPC(tunnelContext)
	go WriteData(tunnelContext)
}
