package session

// Pulling the data from the agents

import "log"

func PollGRPC(tunnelCtx *TunnleContext) {

	dataReceived := false

	for {
		stream := *tunnelCtx.stream
		conn := tunnelCtx.tcp

		msg, err := stream.Recv()

		if err != nil {
			if !dataReceived {
				log.Printf("[PollGRPC] WARNING: stream closed before agent sent any response data connection_id=%s err=%v — TLS handshake may have failed on the local server",
					tunnelCtx.connection_id, err)
			} else {
				log.Printf("[PollGRPC] stream recv error connection_id=%s err=%v", tunnelCtx.connection_id, err)
			}
			conn.Close()
			return
		}

		data := msg.GetData()

		if data == nil || data.ConnectionId != tunnelCtx.connection_id {
			continue
		}

		log.Printf("[PollGRPC] writing %d bytes to client connection_id=%s", len(data.Payload), tunnelCtx.connection_id)
		_, err = conn.Write(data.Payload)
		if err != nil {
			log.Printf("[PollGRPC] write error connection_id=%s err=%v", tunnelCtx.connection_id, err)
			sendClose(tunnelCtx, "Write poll error")
			return
		}
		dataReceived = true
		log.Printf("[PollGRPC] wrote %d bytes to client connection_id=%s", len(data.Payload), tunnelCtx.connection_id)
	}
}
