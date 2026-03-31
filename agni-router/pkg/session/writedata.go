package session

// This function to write the data in the agent

import (
	"time"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
)

func WriteData(tunnelCtx *TunnleContext) {

	stream := *tunnelCtx.stream
	conn := tunnelCtx.tcp
	start := time.Now()
	firstRead := true
	logger.Logger.Info("relay write loop started", "connection_id", tunnelCtx.connection_id)

	buf := make([]byte, 32*1024)

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			if firstRead {
				logger.Logger.Warn("connection closed before any data received",
					"connection_id", tunnelCtx.connection_id,
					"duration_ms", time.Since(start).Milliseconds(),
					"error", err,
					"hint", "client may have rejected the TLS certificate",
				)
			} else {
				logger.Logger.Error("read error from client", "connection_id", tunnelCtx.connection_id, "error", err)
			}
			sendClose(tunnelCtx, "buffer read error")
			return
		}
		firstRead = false
		logger.Logger.Info("read bytes from client", "connection_id", tunnelCtx.connection_id, "bytes", n)

		err = stream.Send(&tunnel.Envelope{
			Message: &tunnel.Envelope_Data{
				Data: &tunnel.TunnelData{
					Payload:      append([]byte(nil), buf[:n]...),
					ConnectionId: tunnelCtx.connection_id,
				},
			},
		})

		if err != nil {
			logger.Logger.Error("failed to send payload to agent", "connection_id", tunnelCtx.connection_id, "error", err)
			conn.Close()
			return
		}
	}
}
