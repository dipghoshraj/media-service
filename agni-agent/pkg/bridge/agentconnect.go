package bridge

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	mp "github.com/odio4u/mem-sdk/memsdk/maps"

	"github.com/odio4u/mem-sdk/memsdk/pkg"
)

func AgentFingerprint() (string, error) {
	permfile := fmt.Sprintf("%s/client.pem", YamlConfig.Agent.Certs)
	certPEM, err := os.ReadFile(permfile)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(certPEM)

	if block == nil || block.Type != "CERTIFICATE" {
		return "", err
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(cert.Raw)
	fingerprint := hex.EncodeToString(sum[:])
	Logger.Info("computed agent cert fingerprint", "fingerprint", fingerprint)
	return fingerprint, nil
}

func AgentRegistry() (mp.Agent, string, error) {

	config := pkg.Config{
		Address:     YamlConfig.Agent.Seeder.Address,
		Fingerprint: YamlConfig.Agent.Seeder.Fingureprint,
		Timeout:     5 * time.Second,
	}

	client, err := mp.NewSdkOperation(config)
	if err != nil {
		return mp.Agent{}, "", err
	}

	Logger.Info("looking up gateways", "region", YamlConfig.Agent.Region)

	gateways, err := client.GetGatewayInfo(context.Background(), YamlConfig.Agent.Region)
	if err != nil {
		return mp.Agent{}, "", err
	}

	if len(gateways) == 0 {
		return mp.Agent{}, "", fmt.Errorf("no gateways found")
	}

	gw := gateways[0]
	Logger.Info("resolved gateway", "gateway_id", gw.ID, "gateway_ip", gw.IP)

	fingerprint, err := AgentFingerprint()
	if err != nil {
		return mp.Agent{}, "", err
	}
	Logger.Info("using agent fingerprint", "fingerprint", fingerprint)

	agent, err := client.ConnectAgent(context.Background(), YamlConfig.Agent.Domain, gw.ID, fingerprint, YamlConfig.Agent.Region)
	if err != nil {
		return mp.Agent{}, "", err
	}

	Logger.Info("agent registered", "agent_id", agent.ID, "domain", agent.Domain)
	return agent, gw.Identity, nil
}
