package rpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/odio4u/agni-schema/tunnel"
)

func (ts *TunnelSession) HandleStream(ctx context.Context, connectionid string, payload []byte) error {

	localconn := ts.Localconn
	conn, exist := localconn.LocalConn[connectionid]
	if !exist {
		return fmt.Errorf("Falied to fetch the connection from fabric")
	}
	_, err := conn.Write(payload)
	if err != nil {
		return fmt.Errorf("Can not write to local server %s", connectionid)
	}
	log.Println("[Agni Agent] write data connection id: and size:", connectionid, len(payload))
	return nil
}

func (ts *TunnelSession) LocaltoRpc(ctx context.Context, conn net.Conn, connectionid string) {
	go func() {
		buf := make([]byte, 32*1024)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					log.Println("Local connection closed:", connectionid)
				} else {
					log.Println("Local read error:", err)
				}
				return
			}

			if n == 0 {
				continue
			}

			payload := append([]byte(nil), buf[:n]...)

			log.Println("[Agni Agent] Write data to router : size :", connectionid, len(payload))
			ts.sendMu.Lock()
			err = ts.Stream.Send(&tunnel.Envelope{
				Message: &tunnel.Envelope_Data{
					Data: &tunnel.TunnelData{
						ConnectionId: connectionid,
						Payload:      payload,
					},
				},
			})
			ts.sendMu.Unlock()

			if err != nil {
				log.Println("[Agni Agent] send failed:", err)
				return
			}
		}
	}()
}
