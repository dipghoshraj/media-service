package connects

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

func (c *TunnelClient) Start(ctx context.Context) error {
	var dialer websocket.Dialer

	conn, _, err := dialer.DialContext(ctx, c.cfg.GatewayURL, nil)
	if err != nil {
		log.Printf("Failed to connect to gateway: %v", err)
		return err
	}

	c.conn = conn
	log.Printf("Connected to gateway: %s", c.cfg.GatewayURL)

	// Need to be logic for the long lived connection and handshake process

	close(c.close)
	return nil

}
