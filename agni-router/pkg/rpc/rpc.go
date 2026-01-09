package rpc

// github.com/Purple-House/agni-tunnels/tunnel-proto

import (
	tunnel "github.com/Purple-House/agni-schema/protobuf"
	mp "github.com/Purple-House/mem-sdk/memsdk/maps"
)

type TunnelRpc struct {
	tunnel.UnimplementedAgniTunnelServer
}

func NewTunnelRpc(registry *mp.Client) *TunnelRpc {
	return &TunnelRpc{}
}
