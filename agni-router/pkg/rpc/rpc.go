package rpc

// github.com/odio4u/agni-tunnels/tunnel-proto

import (
	tunnel "github.com/odio4u/agni-schema/tunnel"
)

type TunnelRpc struct {
	tunnel.UnimplementedAgniTunnelServer
}
