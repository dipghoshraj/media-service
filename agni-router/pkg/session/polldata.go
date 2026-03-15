package session

// Pulling the data from the agents

import "github.com/odio4u/agni-tunnels/agni-router/pkg/logger"

func PollGRPC(tunnelCtx *TunnleContext) {

	dataReceived := false

	for {
		stream := *tunnelCtx.stream
		conn := tunnelCtx.tcp

		msg, err := stream.Recv()

		if err != nil {
			if !dataReceived {
				logger.Logger.Warn("stream closed before agent sent any response",
					"connection_id", tunnelCtx.connection_id,
					"error", err,
					"hint", "TLS handshake may have failed on the local server",
				)
			} else {
				logger.Logger.Error("stream recv error", "connection_id", tunnelCtx.connection_id, "error", err)
			}
			conn.Close()
			return
		}

		data := msg.GetData()

		if data == nil || data.ConnectionId != tunnelCtx.connection_id {
			continue
		}

		logger.Logger.Info("writing bytes to client", "connection_id", tunnelCtx.connection_id, "bytes", len(data.Payload))
		_, err = conn.Write(data.Payload)
		if err != nil {
			logger.Logger.Error("write to client failed", "connection_id", tunnelCtx.connection_id, "error", err)
			sendClose(tunnelCtx, "Write poll error")
			return
		}
		dataReceived = true
		logger.Logger.Info("wrote bytes to client", "connection_id", tunnelCtx.connection_id, "bytes", len(data.Payload))
	}
}
