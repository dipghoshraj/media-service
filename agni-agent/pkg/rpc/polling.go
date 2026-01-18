package rpc

import (
	"log"

	tunnel "github.com/odio4u/agni-schema/tunnel"
)

func PollStream(stream tunnel.AgniTunnel_ConnectClient) {

	for {

		msg, err := stream.Recv()
		if err != nil {
			log.Println("[Agni-Angent] Polling error stream closed ", err)
		}

		switch m := msg.Message.(type) {
		case *tunnel.Envelope_ConnectAck:
			log.Println("[Agni-Agent] Connection Ack", m.ConnectAck.Accepted)
		case *tunnel.Envelope_Close:
			log.Println("[Agni-Agent] Connection closed", m.Close.ConnectionId)
		case *tunnel.Envelope_Open:
			log.Println("[Agni-Agent] Connection open", m.Open.ConnectionId)
		case *tunnel.Envelope_Data:
			log.Println("[Agni-Agent] Connection Dara", m.Data.ConnectionId)
		default:
			log.Printf("[Agni-Agent] Unkown event type %T", m)
		}
	}
}
