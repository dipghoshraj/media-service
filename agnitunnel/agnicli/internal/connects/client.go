package connects

import (
	"agni-cli/internal/types"
	"crypto/rand"
	"encoding/hex"
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

func GenerateNonce(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func NewTunnelClient(cfg ClientConfig) (*TunnelClient, error) {
	return &TunnelClient{
		Cfg:   cfg,
		Close: make(chan struct{}),
	}, nil
}

// log.Printf("Agent [%s] write error: %v", sess.AppID, err)
