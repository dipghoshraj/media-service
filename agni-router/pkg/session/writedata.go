package session

import tunnel "github.com/Purple-House/agni-schema/protobuf"

func WriteData(tunnelCtx *TunnleContext) {

	stream := *tunnelCtx.stream
	conn := tunnelCtx.tcp

	buf := make([]byte, 32*1024)

	for {

		n, err := conn.Read(buf)
		if err != nil {
			sendClose(tunnelCtx, "buffer read error")
			return
		}

		err = stream.Send(&tunnel.Envelope{
			Message: &tunnel.Envelope_Data{
				Data: &tunnel.TunnelData{
					Payload:      append([]byte(nil), buf[:n]...),
					ConnectionId: tunnelCtx.connection_id,
				},
			},
		})

		if err != nil {
			conn.Close()
			return
		}
	}
}
