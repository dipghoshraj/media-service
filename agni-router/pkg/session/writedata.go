package session

import (
	"log"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
)

func WriteData(tunnelCtx *TunnleContext) {

	stream := *tunnelCtx.stream
	conn := tunnelCtx.tcp
	log.Printf("[WriteData] started connection_id=%s", tunnelCtx.connection_id)

	buf := make([]byte, 32*1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("[WriteData] read error connection_id=%s err=%v", tunnelCtx.connection_id, err)
			sendClose(tunnelCtx, "buffer read error")
			return
		}
		log.Printf("[WriteData] read %d bytes connection_id=%s", n, tunnelCtx.connection_id)

		err = stream.Send(&tunnel.Envelope{
			Message: &tunnel.Envelope_Data{
				Data: &tunnel.TunnelData{
					Payload:      append([]byte(nil), buf[:n]...),
					ConnectionId: tunnelCtx.connection_id,
				},
			},
		})

		if err != nil {
			log.Printf("[WriteData] read error connection_id=%s err=%v", tunnelCtx.connection_id, err)
			conn.Close()
			return
		}
	}
}
