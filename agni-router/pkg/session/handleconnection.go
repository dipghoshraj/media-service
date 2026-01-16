package session

import (
	"net"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
	"github.com/Purple-House/agni-tunnels/agni-router/pkg/session/sni"
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
		return
	}

	registry, exists := Seeder.GetSession(serverName)
	if !exists {
		return
	}

	_ = registry.stream
}
