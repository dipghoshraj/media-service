package session

import (
	"log"
	"sync"

	tunnel "github.com/odio4u/agni-schema/tunnel"
)

type AgentSession struct {
	AppID  string
	Stream *tunnel.AgniTunnel_ConnectServer
}

type AgentSeeder struct {
	sync.RWMutex
	sessions  map[string]*AgentSession
	domainmap map[string]string
}

var Seeder = &AgentSeeder{
	sessions:  make(map[string]*AgentSession),
	domainmap: make(map[string]string),
}

func (r *AgentSeeder) AddDomainMap(appID string, domain string) {
	r.Lock()
	defer r.Unlock()
	r.domainmap[appID] = domain
	log.Printf("[Agni Router] mapped app: %s with domain %s", appID, domain)
}

func (r *AgentSeeder) GetDomainMap(appID string) (string, bool) {
	r.RLock()
	defer r.RUnlock()
	domain, exist := r.domainmap[appID]
	return domain, exist
}

func (r *AgentSeeder) Register(appID string, session *AgentSession) {
	r.Lock()
	defer r.Unlock()
	r.sessions[appID] = session
	log.Printf("Agent [%s] registered", appID)

}

func (r *AgentSeeder) Unregister(appID string) {
	r.Lock()
	defer r.Unlock()
	if _, exists := r.sessions[appID]; exists {
		delete(r.sessions, appID)
		log.Printf("Agent [%s] unregistered", appID)
	} else {
		log.Printf("Attempted to unregister non-existent agent [%s]", appID)
	}
}

func (r *AgentSeeder) GetSession(appID string) (*AgentSession, bool) {
	r.RLock()
	defer r.RUnlock()
	session, exists := r.sessions[appID]
	if !exists {
		log.Printf("Session for agent [%s] not found", appID)
		return nil, false
	}
	return session, true
}
