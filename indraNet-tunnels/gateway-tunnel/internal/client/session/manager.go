package session

import (
	"log"
	"sync"
	"time"

	pb "gateway-tunnel/proto"

	"github.com/gorilla/websocket"
)

// Gateway tunnel service is a statefull service and this
// agent session details are the most important data service needs to maintain.
type AgentSession struct {
	AppID    string
	Conn     *websocket.Conn
	SendChan chan *pb.TunnelRequest
	LastSeen time.Time
}

type agentRegistry struct {
	sync.RWMutex
	sessions map[string]*AgentSession
}

var Registry = &agentRegistry{
	sessions: make(map[string]*AgentSession),
}

func (r *agentRegistry) Register(appID string, session *AgentSession) {
	r.Lock()
	defer r.Unlock()
	r.sessions[appID] = session
	log.Printf("Agent [%s] registered", appID)

}

func (r *agentRegistry) Unregister(appID string) {
	r.Lock()
	defer r.Unlock()
	if _, exists := r.sessions[appID]; exists {
		delete(r.sessions, appID)
		log.Printf("Agent [%s] unregistered", appID)
	} else {
		log.Printf("Attempted to unregister non-existent agent [%s]", appID)
	}
}

func (r *agentRegistry) GetSession(appID string) (*AgentSession, bool) {
	r.RLock()
	defer r.RUnlock()
	session, exists := r.sessions[appID]
	if !exists {
		log.Printf("Session for agent [%s] not found", appID)
		return nil, false
	}
	return session, true
}

func (r *agentRegistry) Is_online(appID string) bool {
	r.RLock()
	defer r.RUnlock()
	session, exists := r.sessions[appID]
	if !exists {
		log.Printf("Session for agent [%s] not found", appID)
		return false
	}
	// Check if the session is still active based on LastSeen timestamp
	if time.Since(session.LastSeen) > 5*time.Minute {
		log.Printf("Session for agent [%s] is offline", appID)
		return false
	}
	return true
}
