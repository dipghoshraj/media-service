package rpc

import (
	"context"
	"fmt"
	"log"

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
	return nil
}

func (ts *TunnelSession) LocaltoRpc(ctx context.Context, connectionid string) {
	localconn := ts.Localconn
	conn, exist := localconn.LocalConn[connectionid]
	if !exist {
		log.Println("[Agni Agent] Falied to fetch the connection from fabric")
	}

	go func() {
		buf := make([]byte, 32*1024)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				// TODO : implemet close logic
			}

			err = ts.Stream.Send(
				&tunnel.Envelope{
					Message: &tunnel.Envelope_Data{
						Data: &tunnel.TunnelData{
							ConnectionId: connectionid,
							Payload:      append([]byte(nil), buf[:n]...),
						},
					},
				},
			)

			if err != nil {
				// ctx.sendClose("grpc_send_failed")
				return
			}
		}
	}()
}
