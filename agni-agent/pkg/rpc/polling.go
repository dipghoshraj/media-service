package rpc

import (
	"context"
	"io"
	"log"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/connector"
)

func (ts *TunnelSession) PollStream(ctx context.Context) error {

	for {

		select {
		case <-ctx.Done():
			log.Println("[Agni-Agent] PollStream stopped:", ctx.Err())
			return ctx.Err()
		default:
			msg, err := ts.Stream.Recv()
			if err != nil {
				if err == io.EOF {
					log.Println("[Agni-Agent] Stream closed by server")
					return nil
				}
				log.Println("[Agni-Agent] Stream recv error:", err)
				return err
			}
			ts.handleMessage(ctx, msg)
		}

	}
}

func (ts *TunnelSession) handleMessage(ctx context.Context, msg *tunnel.Envelope) {
	switch m := msg.Message.(type) {

	case *tunnel.Envelope_ConnectAck:
		log.Println("[Agni-Agent] Connection Ack:", m.ConnectAck.Accepted)

	case *tunnel.Envelope_Open:
		ts.mu.Lock()
		defer ts.mu.Unlock()

		connection_id := m.Open.ConnectionId
		conn, err := connector.BuildConn(ctx, connection_id)
		if err != nil {
			// TODO : Send connection close
		}

		log.Println("[Agni Agent] connected to local server")
		ts.Localconn = *conn

		go ts.LocaltoRpc(ctx, ts.Localconn.LocalConn[connection_id], connection_id)
		log.Println("[Agni-Agent] Connection open:", connection_id)

	case *tunnel.Envelope_Data:
		err := ts.HandleStream(ctx, m.Data.ConnectionId, m.Data.Payload)
		if err != nil {
			log.Println("[Agni Agent] failed stream:  ", err.Error())
		}
		log.Println("[Agni-Agent] Connection data:", m.Data.ConnectionId)

	case *tunnel.Envelope_Close:
		connection_id := m.Close.ConnectionId
		log.Println("[Agni-Agent] Connection closed:", connection_id)

	default:
		log.Printf("[Agni-Agent] Unknown event type: %T", m)
	}
}
