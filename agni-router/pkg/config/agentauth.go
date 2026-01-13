package config

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"log"
)

func AuthAgent(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

	if len(rawCerts) == 0 {
		return errors.New("no client certificate provided")
	}
	log.Println("coming with client certificates")

	clientCert, err := x509.ParseCertificate(rawCerts[0])
	if err != nil {
		return err
	}

	fp := sha256.Sum256(clientCert.Raw)
	log.Println("Client fingerprint:", hex.EncodeToString(fp[:]))

	if len(clientCert.DNSNames) == 0 {
		return errors.New("no DNS SAN present in client cert")
	}
	agentID := clientCert.DNSNames[0]

	log.Println("Authenticating agent:", agentID)

	identity, err := getAgent(agentID)
	if err != nil {
		log.Println("can not fetch the identity of the agent")
		return errors.New("No agent found")
	}

	log.Printf("%s agent with idenity %s", agentID, identity)

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

	return agentdata.Identity, nil
}
