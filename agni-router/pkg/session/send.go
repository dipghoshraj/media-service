package session

import (
	tunnel "github.com/odio4u/agni-schema/tunnel"
)

func sendClose(tctx *TunnleContext, reason string) {
	tctx.closeOnce.Do(func() {
		stream := *tctx.stream
		_ = stream.Send(
			&tunnel.Envelope{
				Message: &tunnel.Envelope_Close{
					Close: &tunnel.TunnelClose{
						ConnectionId: tctx.connection_id,
						Reason:       reason,
					},
				},
			},
		)

		tctx.tcp.Close()
		close(tctx.closed)
	})
}

func sendOpen(tctx *TunnleContext) {
	stream := *tctx.stream
	_ = stream.Send(
		&tunnel.Envelope{
			Message: &tunnel.Envelope_Open{
				Open: &tunnel.TunnelOpen{
					ConnectionId: tctx.connection_id,
				},
			},
		},
	)
}
