package session

import (
	"net"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
	"github.com/Purple-House/agni-tunnels/agni-router/pkg/session/sni"
	"github.com/google/uuid"
)

type TunnleContext struct {
	connection_id string
	stream        *tunnel.AgniTunnel_ConnectServer
	tcp           net.Conn
	closed        chan struct{}
}

func HandleStream(conn net.Conn) {
	// session.AgentRegistry.GetSession()

	serverName, _, err := sni.PeekSNI(conn)
	if err != nil {
		conn.Close()
		return
	}

	session, exists := Seeder.GetSession(serverName)
	if !exists {
		conn.Close()
		return
	}
	tunnelContext := &TunnleContext{
		connection_id: uuid.New().String(),
		stream:        session.stream,
		tcp:           conn,
		closed:        make(chan struct{}),
	}

	go PollGRPC(tunnelContext)
	go WriteData(tunnelContext)
}
