package session

import (
	"log"
	"net"

	"github.com/google/uuid"
	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/mem-sdk/sni"
)

type TunnleContext struct {
	connection_id string
	stream        *tunnel.AgniTunnel_ConnectServer
	tcp           net.Conn
	closed        chan struct{}
}

func HandleStream(conn net.Conn) {
	// session.AgentRegistry.GetSession()
	log.Println("Coming on the handle stream")

	serverName, wrappedConn, err := sni.PeekSNI(conn)
	if err != nil {
		conn.Close()
		return
	}
	log.Println("Received the connection from the proxy", serverName)

	session, exists := Seeder.GetSession(serverName)
	if !exists {
		conn.Close()
		return
	}

	tunnelContext := &TunnleContext{
		connection_id: uuid.New().String(),
		stream:        session.Stream,
		tcp:           wrappedConn,
		closed:        make(chan struct{}),
	}

	sendOpen(tunnelContext)

	go PollGRPC(tunnelContext)
	go WriteData(tunnelContext)
}
