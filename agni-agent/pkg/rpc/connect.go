package rpc

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/connector"
	"github.com/odio4u/mem-sdk/memsdk/maps"
	"google.golang.org/grpc"
)

type TunnelSession struct {
	Conn      *grpc.ClientConn
	Stream    tunnel.AgniTunnel_ConnectClient
	Ctx       context.Context
	Cancel    context.CancelFunc
	Close     chan struct{}
	Localconn connector.LocalConn
	mu        sync.Mutex
	sendMu    sync.Mutex
}

func InitateConnection(router string, gatewayIdentity string) *grpc.ClientConn {
	conn := routerConnect(router, gatewayIdentity)
	return conn
}

// build the local server connection
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
				Token:     agent.Domain,
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

// Send connection is the function which will start the transaction
// with the local server

func (ts *TunnelSession) SendConnection() {

	done := make(chan struct{})
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := ts.PollStream(ctx); err != nil {
			log.Println("[Agni-Agent] PollStream exited:", err)
		}
		close(done)
	}()

	<-quit
	log.Println("Shutting down connection...")
	ts.Cancel()
	ts.Conn.Close()
	<-done
}
