package session

import (
	"sync"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
)

func sendClose(tctx *TunnleContext, reason string) {
	var once sync.Once

	once.Do(func() {
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
