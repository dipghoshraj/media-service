package rpc

// github.com/Purple-House/agni-tunnels/tunnel-proto

import (
	tunnel "github.com/Purple-House/agni-schema/protobuf"
)

type TunnelRpc struct {
	tunnel.UnimplementedAgniTunnelServer
}
