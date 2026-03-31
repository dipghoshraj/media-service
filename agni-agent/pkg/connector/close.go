package connector

import (
	"context"
	"sync"

	"github.com/odio4u/agni-schema/tunnel"
)

func SendClose(ctx context.Context, reason string, stream tunnel.AgniTunnel_ConnectServer) error {

	var once sync.Once
	once.Do(func() {
		connection_id := ctx.Value("connection_id").(string)
		_ = stream.Send(
			&tunnel.Envelope{
				Message: &tunnel.Envelope_Close{
					Close: &tunnel.TunnelClose{
						ConnectionId: connection_id,
						Reason:       reason,
					},
				},
			},
		)
	})

	return nil
}
