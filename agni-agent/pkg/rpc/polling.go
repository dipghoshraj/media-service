package rpc

import (
	"context"
	"fmt"
	"io"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/connector"
)

func (ts *TunnelSession) PollStream(ctx context.Context) error {

	for {
		select {
		case <-ctx.Done():
			bridge.Logger.Info("poll stream stopped", "reason", ctx.Err())
			return ctx.Err()
		default:
			msg, err := ts.Stream.Recv()
			if err != nil {
				if err == io.EOF {
					bridge.Logger.Info("stream closed by server")
					return nil
				}
				bridge.Logger.Error("stream recv error", "error", err)
				return err
			}
			ts.handleMessage(ctx, msg)
		}
	}
}

func (ts *TunnelSession) handleMessage(ctx context.Context, msg *tunnel.Envelope) {
	switch m := msg.Message.(type) {

	case *tunnel.Envelope_ConnectAck:
		bridge.Logger.Info("connect ack received", "accepted", m.ConnectAck.Accepted)

	case *tunnel.Envelope_Open:
		ts.mu.Lock()
		defer ts.mu.Unlock()

		connection_id := m.Open.ConnectionId
		conn, err := connector.BuildConn(ctx, connection_id)
		if err != nil {
			bridge.Logger.Warn("failed to dial local server",
				"connection_id", connection_id,
				"error", err,
				"hint", "is the local app running on the configured forward port?",
			)
			return
		}

		bridge.Logger.Info("connected to local server", "connection_id", connection_id)
		ts.Localconn = *conn

		go ts.LocaltoRpc(ctx, ts.Localconn.LocalConn[connection_id], connection_id)
		bridge.Logger.Info("tunnel open", "connection_id", connection_id)

	case *tunnel.Envelope_Data:
		err := ts.HandleStream(ctx, m.Data.ConnectionId, m.Data.Payload)
		if err != nil {
			bridge.Logger.Error("failed to handle stream data", "connection_id", m.Data.ConnectionId, "error", err)
		}
		bridge.Logger.Info("data frame received", "connection_id", m.Data.ConnectionId)

	case *tunnel.Envelope_Close:
		connection_id := m.Close.ConnectionId
		bridge.Logger.Info("tunnel closed", "connection_id", connection_id)

	default:
		bridge.Logger.Warn("unknown envelope type", "type", fmt.Sprintf("%T", m))
	}
}
