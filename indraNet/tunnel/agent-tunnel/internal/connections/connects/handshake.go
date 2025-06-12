package connects

import (
	"agent-tunnel/proto"
	"context"
	"fmt"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

func (c *TunnelClient) Handshake(ctx context.Context) error {
	nonce := generateNonce()
	timestamp := time.Now().Unix()

	msg := fmt.Sprintf("%s:%d:%s", c.Cfg.Token, timestamp, nonce)
	signature := c.signature(msg)

	connectReq := &proto.ConnectRequest{
		AgentId:   c.Cfg.AgentID,
		Token:     c.Cfg.Token,
		Timestamp: timestamp,
		Nonce:     nonce,
		Signature: signature,
	}

	env := &proto.Envelope{
		Message: &proto.Envelope_Connect{
			Connect: connectReq,
		},
	}

	data, err := protobuf.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal connect request: %w", err)
	}

	err = c.conn.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		return fmt.Errorf("write connect message: %w", err)
	}
	return nil
}
