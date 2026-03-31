package config

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"

	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/session"
)

func AuthAgent(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

	if len(rawCerts) == 0 {
		return errors.New("no client certificate provided")
	}
	logger.Logger.Info("client certificate received")

	clientCert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return err
	}

	fp := sha256.Sum256(clientCert.Raw)

	if len(clientCert.DNSNames) == 0 {
		return errors.New("no DNS SAN present in client cert")
	}
	agentID := clientCert.DNSNames[0]

	logger.Logger.Info("authenticating agent", "agent_id", agentID, "fingerprint", hex.EncodeToString(fp[:]))

	identity, err := getAgent(agentID)
	if err != nil {
		logger.Logger.Error("failed to fetch agent identity", "agent_id", agentID, "error", err)
		return errors.New("No agent found")
	}

	logger.Logger.Info("agent identity resolved", "agent_id", agentID, "identity", identity)

	if hex.EncodeToString(fp[:]) != identity {
		return errors.New("client fingerprint mismatch")
	}
	return nil
}

func getAgent(agentID string) (string, error) {

	seedClient, err := SeederClient()
	if err != nil {
		return "", err
	}

	ctx := context.Context(context.Background())

	agentdata, err := seedClient.GetAgentProxyMapping(ctx, YamlConfig.Router.Region, agentID)
	if err != nil {
		return "", err
	}
	session.Seeder.AddDomainMap(agentdata.ID, agentdata.Domain)
	return agentdata.Identity, nil
}
