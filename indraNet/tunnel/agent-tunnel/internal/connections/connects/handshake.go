package connects

import (
	"agent-tunnel/proto"
	"context"
	"fmt"
	"time"
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

	_ = c.send_envalope(env)
	return nil
}
