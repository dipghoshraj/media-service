package rpc

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
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
	bridge.Logger.Info("wrote data to local server", "connection_id", connectionid, "bytes", len(payload))
	return nil
}

func (ts *TunnelSession) LocaltoRpc(ctx context.Context, conn net.Conn, connectionid string) {
	go func() {
		buf := make([]byte, 32*1024)

		for {
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					bridge.Logger.Info("local connection closed", "connection_id", connectionid)
				} else {
					bridge.Logger.Error("local read error", "connection_id", connectionid, "error", err)
				}
				ts.sendMu.Lock()
				_ = ts.Stream.Send(&tunnel.Envelope{
					Message: &tunnel.Envelope_Close{
						Close: &tunnel.TunnelClose{
							ConnectionId: connectionid,
							Reason:       "remote connection closed",
						},
					},
				})
				ts.sendMu.Unlock()
				return
			}

			if n == 0 {
				bridge.Logger.Info("local aggregator returned 0 bytes", "connection_id", connectionid)
				continue
			}

			bridge.Logger.Info("received bytes from local app", "connection_id", connectionid, "bytes", n)
			payload := append([]byte(nil), buf[:n]...)

			bridge.Logger.Info("forwarding payload to router", "connection_id", connectionid, "bytes", len(payload))
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
				bridge.Logger.Error("failed to send payload to router", "connection_id", connectionid, "error", err)
				return
			}
		}
	}()
}
