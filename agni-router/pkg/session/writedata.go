package session

// This function to write the data in the agent

import (
	"log"
	"time"

	tunnel "github.com/odio4u/agni-schema/tunnel"
)

func WriteData(tunnelCtx *TunnleContext) {

	stream := *tunnelCtx.stream
	conn := tunnelCtx.tcp
	start := time.Now()
	firstRead := true
	log.Printf("[WriteData] started connection_id=%s", tunnelCtx.connection_id)

	buf := make([]byte, 32*1024)

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			if firstRead {
				log.Printf("[WriteData] WARNING: connection closed before any data was received connection_id=%s duration=%v err=%v — client may have rejected the TLS certificate",
					tunnelCtx.connection_id, time.Since(start).Round(time.Millisecond), err)
			} else {
				log.Printf("[WriteData] read error connection_id=%s err=%v", tunnelCtx.connection_id, err)
			}
			sendClose(tunnelCtx, "buffer read error")
			return
		}
		firstRead = false
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
