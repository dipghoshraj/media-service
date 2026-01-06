package bridge

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	mp "github.com/Purple-House/mem-sdk/memsdk/maps"

	"github.com/Purple-House/mem-sdk/memsdk/pkg"
)

func AgentFingerprint() (string, error) {
	permfile := "client.pem"
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
	log.Printf("Client CERT fingerprint (SHA256): %s", fingerprint)
	return fingerprint, nil
}

func AgentRegistry() (mp.Agent, string, error) {

	config := pkg.Config{
		Address:     YamlConfig.Agent.Registry.Address,
		Fingerprint: YamlConfig.Agent.Registry.Fingureprint,
		Timeout:     5 * time.Second,
	}

	client, err := mp.NewSdkOperation(config)
	if err != nil {
		return mp.Agent{}, "", err
	}

	log.Printf("finding the gateways in %s region", YamlConfig.Agent.Region)

	gateways, err := client.GetGatewayInfo(context.Background(), YamlConfig.Agent.Region)
	if err != nil {
		return mp.Agent{}, "", err
	}

	if len(gateways) == 0 {
		return mp.Agent{}, "", fmt.Errorf("no gateways found")
	}

	gw := gateways[0]
	log.Printf("Resolved Gateway: %s (%s)\n", gw.IP, gw.ID)

	fingerprint, err := AgentFingerprint()
	if err != nil {
		return mp.Agent{}, "", err
	}
	log.Printf("Using Agent Fingerprint: %s\n", fingerprint)

	agent, err := client.ConnectAgent(context.Background(), YamlConfig.Agent.Domain, gw.ID, fingerprint, YamlConfig.Agent.Region)
	if err != nil {
		return mp.Agent{}, "", err
	}

	log.Printf("Connected Agent: %s (%s)\n", agent.ID, agent.Domain)
	return agent, fingerprint, nil
}
