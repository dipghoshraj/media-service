package session

import (
	"log"
	"sync"
	"time"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
)

type AgentSession struct {
	AppID    string
	stream   *tunnel.AgniTunnel_ConnectServer
	SendChan chan *tunnel.TunnelOpen
	LastSeen time.Time
}

type AgentRegistry struct {
	sync.RWMutex
	sessions map[string]*AgentSession
}

var Registry = &AgentRegistry{
	sessions: make(map[string]*AgentSession),
}

func (r *AgentRegistry) Register(appID string, session *AgentSession) {
	r.Lock()
	defer r.Unlock()
	r.sessions[appID] = session
	log.Printf("Agent [%s] registered", appID)

}

func (r *AgentRegistry) Unregister(appID string) {
	r.Lock()
	defer r.Unlock()
	if _, exists := r.sessions[appID]; exists {
		delete(r.sessions, appID)
		log.Printf("Agent [%s] unregistered", appID)
	} else {
		log.Printf("Attempted to unregister non-existent agent [%s]", appID)
	}
}

func (r *AgentRegistry) GetSession(appID string) (*AgentSession, bool) {
	r.RLock()
	defer r.RUnlock()
	session, exists := r.sessions[appID]
	if !exists {
		log.Printf("Session for agent [%s] not found", appID)
		return nil, false
	}
	return session, true
}
