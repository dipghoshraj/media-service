package session

import (
	"sync"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
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
	logger.Logger.Info("domain mapping added", "agent_id", appID, "domain", domain)
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
	logger.Logger.Info("agent session registered", "agent_id", appID)
}

func (r *AgentSeeder) Unregister(appID string) {
	r.Lock()
	defer r.Unlock()
	if _, exists := r.sessions[appID]; exists {
		delete(r.sessions, appID)
		logger.Logger.Info("agent session unregistered", "agent_id", appID)
	} else {
		logger.Logger.Warn("unregister called for unknown agent", "agent_id", appID)
	}
}

func (r *AgentSeeder) GetSession(appID string) (*AgentSession, bool) {
	r.RLock()
	defer r.RUnlock()
	session, exists := r.sessions[appID]
	if !exists {
		logger.Logger.Warn("session not found", "agent_id", appID)
		return nil, false
	}
	return session, true
}
