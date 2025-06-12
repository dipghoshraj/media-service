package connects

import (
	"agent-tunnel/internal/types"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ClientConfig = types.ClientConfig

type TunnelClient struct {
	cfg   ClientConfig
	conn  *websocket.Conn
	close chan struct{}

	mu      sync.Mutex
	streams map[string]net.Conn // Stream ID → local TCP conn
}

func generateNonce() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func NewTunnelClient(cfg ClientConfig) (*TunnelClient, error) {
	return &TunnelClient{
		cfg:   cfg,
		close: make(chan struct{}),
	}, nil
}
