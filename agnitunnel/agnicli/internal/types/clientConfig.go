package types

type ClientConfig struct {
	GatewayURL  string
	Token       string
	Secret      string
	AgentID     string
	Portforward string // Local port to forward

}
