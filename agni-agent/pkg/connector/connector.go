package connector

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
)

type LocalConn struct {
	LocalConn map[string]net.Conn
}

func BuildConn(ctx context.Context, id string) (*LocalConn, error) {
	ctx = context.WithValue(ctx, "connection_id", id)

	localconn, err := NewConnectionCtx()
	if err != nil {
		return nil, err
	}

	connmap := make(map[string]net.Conn)
	connmap[id] = localconn
	return &LocalConn{
		LocalConn: connmap,
	}, nil
}

func NewConnectionCtx() (net.Conn, error) {
	port := bridge.YamlConfig.Agent.Forward
	host := bridge.YamlConfig.Agent.Host

	url := net.JoinHostPort(host, fmt.Sprintf("%d", port))

	localConn, err := net.Dial("tcp", url)
	if err != nil {
		return nil, errors.New("Failed to connect with local connection pool")
	}

	return localConn, nil
}
