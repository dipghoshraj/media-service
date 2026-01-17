package rpc

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/mem-sdk/memsdk/maps"
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
		session.Cancel()
		session.Conn.Close()
		panic(err)
	}

	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		ReadLoop(session.Stream)
		close(done)
	}()

	<-quit
	log.Println("Shutting down connection...")
	session.Cancel()
	session.Conn.Close()
	<-done
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
