package connects

import (
	tunnelv1 "agni-cli/proto"
	"log"
)

func (t *TunnelClient) HandleMessages(env *tunnelv1.Envelope) error {
	switch m := env.Message.(type) {
	case *tunnelv1.Envelope_Open:
		log.Printf("Opening tunnel: %s\n", m.Open.Id)
		// Handle opening a tunnel
	}
	return nil
}
