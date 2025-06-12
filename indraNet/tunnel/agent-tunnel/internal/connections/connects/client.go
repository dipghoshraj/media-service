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
	Cfg   ClientConfig
	conn  *websocket.Conn
	Close chan struct{}

	mu      sync.Mutex
	Streams map[string]net.Conn // Stream ID → local TCP conn
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
