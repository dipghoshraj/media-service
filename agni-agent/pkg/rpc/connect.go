package rpc

import (
	"context"
	"log"
	"time"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
	"github.com/Purple-House/mem-sdk/memsdk/maps"
	"google.golang.org/grpc"
)

func InitateConnection(router string) *grpc.ClientConn {
	conn := routerConnect(router)
	return conn
}

func SendConnection(agent maps.Agent, fingerprint string) {

	conn := GetRouter()

	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client := tunnel.NewAgniTunnelClient(conn)

	stream, err := client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = stream.Send(&tunnel.Envelope{
		Message: &tunnel.Envelope_Connect{
			Connect: &tunnel.ConnectRequest{
				AgentId:   agent.ID,
				Token:     fingerprint,
				Timestamp: time.Now().Unix(),
				Nonce:     "",
				Signature: fingerprint,
			},
		},
	})
	if err != nil {
		log.Fatal("send connect:", err)
	}

	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				log.Println("stream closed:", err)
				return
			}

			log.Println("received the message: ", in.Message)

		}
	}()
}
