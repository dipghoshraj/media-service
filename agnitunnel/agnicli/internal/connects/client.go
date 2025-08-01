package connects

import (
	"agni-cli/internal/types"
	"fmt"
	"net"
	"sync"
	"time"

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

func generateNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewTunnelClient(cfg ClientConfig) (*TunnelClient, error) {
	return &TunnelClient{
		Cfg:   cfg,
		Close: make(chan struct{}),
	}, nil
}

// log.Printf("Agent [%s] write error: %v", sess.AppID, err)
