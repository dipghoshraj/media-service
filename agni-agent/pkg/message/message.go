package message

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
)

type Message struct {
	ctx context.Context
}

type ConnectionCtx struct {
	ID        string
	LocalConn net.Conn
}

func NewMessage() *Message {
	return &Message{
		ctx: context.Background(),
	}
}

func NewConnectionCtx(connection_id string) *ConnectionCtx {
	port := bridge.YamlConfig.Agent.Forward
	host := net.JoinHostPort("localhost", fmt.Sprintf("%d", port))

	localConn, err := net.Dial("tcp", host)
	if err != nil {
		log.Panicln("[Agni Agent] paniced to build local connection")
	}

	return &ConnectionCtx{
		ID:        connection_id,
		LocalConn: localConn,
	}
}
