package rpc

import (
	"context"
	"log"
	"time"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
	"github.com/Purple-House/mem-sdk/memsdk/maps"
	"google.golang.org/grpc"
)

type TunnelSession struct {
	Conn   *grpc.ClientConn
	Stream tunnel.AgniTunnel_ConnectClient
	Ctx    context.Context
	Cancel context.CancelFunc
}

func InitateConnection(router string, gatewayIdentity string) *grpc.ClientConn {
	conn := routerConnect(router, gatewayIdentity)
	return conn
}

func NewTunnelSession(agent maps.Agent) (*TunnelSession, error) {
	conn := GetRouter()

	ctx, cancel := context.WithCancel(context.Background())
	client := tunnel.NewAgniTunnelClient(conn)

	stream, err := client.Connect(ctx)
	if err != nil {
		cancel()
		conn.Close()
		return nil, err
	}

	err = stream.Send(&tunnel.Envelope{
		Message: &tunnel.Envelope_Connect{
			Connect: &tunnel.ConnectRequest{
				AgentId:   agent.ID,
				Token:     agent.Identity,
				Timestamp: time.Now().Unix(),
				Signature: agent.Identity,
			},
		},
	})
	if err != nil {
		cancel()
		conn.Close()
		return nil, err
	}

	return &TunnelSession{
		Conn:   conn,
		Stream: stream,
		Ctx:    ctx,
		Cancel: cancel,
	}, nil
}

func SendConnection(agent maps.Agent) {

	session, err := NewTunnelSession(agent)
	if err != nil {
		log.Fatal(err)
	}

	go ReadLoop(session.Stream)
	session.Cancel()
	session.Conn.Close()
	// select {} // 🚨 blocks forever
}

func ReadLoop(stream tunnel.AgniTunnel_ConnectClient) {
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Println("stream closed:", err)
			return
		}
		log.Println("received the message: ", in.Message)
	}
}
