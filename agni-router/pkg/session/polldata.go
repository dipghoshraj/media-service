package session

func PollGRPC(tunnelCtx *TunnleContext) {

	for {
		stream := *tunnelCtx.stream
		conn := tunnelCtx.tcp

		msg, err := stream.Recv()

		if err != nil {
			conn.Close()
			return
		}

		data := msg.GetData()

		if data == nil || data.ConnectionId != tunnelCtx.connection_id {
			continue
		}

		_, err = conn.Write(data.Payload)
		if err != nil {
			sendClose(tunnelCtx, "Write poll error")
			return
		}
	}
}
