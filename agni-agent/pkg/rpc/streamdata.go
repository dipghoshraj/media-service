package rpc

import (
	"context"
	"fmt"
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
