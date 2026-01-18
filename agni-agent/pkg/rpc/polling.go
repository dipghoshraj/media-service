package rpc

import (
	"context"
	"io"
	"log"

	tunnel "github.com/odio4u/agni-schema/tunnel"
)

func PollStream(ctx context.Context, stream tunnel.AgniTunnel_ConnectClient) error {

	for {

		select {
		case <-ctx.Done():
			log.Println("[Agni-Agent] PollStream stopped:", ctx.Err())
			return ctx.Err()
		default:
			msg, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("[Agni-Agent] Stream closed by server")
					return nil
				}
				log.Println("[Agni-Agent] Stream recv error:", err)
				return err
			}
			handleMessage(msg)
		}

	}
}

func handleMessage(msg *tunnel.Envelope) {
	switch m := msg.Message.(type) {

	case *tunnel.Envelope_ConnectAck:
		log.Println("[Agni-Agent] Connection Ack:", m.ConnectAck.Accepted)

	case *tunnel.Envelope_Open:
		log.Println("[Agni-Agent] Connection open:", m.Open.ConnectionId)

	case *tunnel.Envelope_Data:
		log.Println("[Agni-Agent] Connection data:", m.Data.ConnectionId)

	case *tunnel.Envelope_Close:
		log.Println("[Agni-Agent] Connection closed:", m.Close.ConnectionId)

	default:
		log.Printf("[Agni-Agent] Unknown event type: %T", m)
	}
}
