package connects

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (c *TunnelClient) signature(msg string) string {
	// Placeholder for signature generation logic
	mac := hmac.New(sha256.New, []byte(c.Cfg.Secret))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}
