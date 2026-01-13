package session

import (
	"time"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
)

type AgentSession struct {
	AppID    string
	Conn     *tunnel.AgniTunnel_ConnectServer
	SendChan chan *tunnel.TunnelOpen
	LastSeen time.Time
}
