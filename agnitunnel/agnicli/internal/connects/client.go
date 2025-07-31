package connects

import (
	"agni-cli/internal/types"
	"net"
	"sync"

	tunnelv1 "agni-cli/proto"
)

type ClientConfig = types.ClientConfig

type TunnelClient struct {
	Cfg    ClientConfig
	stream tunnelv1.TunnelService_TunnelStreamClient
	Close  chan struct{}

	mu      sync.Mutex
	Streams map[string]net.Conn // Map of stream IDs to net.Conn
}
