package connects

import (
	"agent-tunnel/internal/types"

	"github.com/gorilla/websocket"
)

type ClientConfig = types.ClientConfig

type TunnelClient struct {
	cfg   ClientConfig
	conn  *websocket.Conn
	close chan struct{}
}

// func generateNonce() string {
// 	return fmt.Sprintf("%d", time.Now().UnixNano())
// }

func NewTunnelClient(cfg ClientConfig) (*TunnelClient, error) {
	return &TunnelClient{
		cfg:   cfg,
		close: make(chan struct{}),
	}, nil
}
