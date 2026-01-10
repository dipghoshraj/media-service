package config

import (
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

func CertFingurePrint() (string, error) {
	permfile := fmt.Sprintf("%s/server.pem", YamlConfig.Router.Certs)
	certPEM, err := os.ReadFile(permfile)
	if err != nil {
		return "", fmt.Errorf("failed to read certificate file: %v", err)
	}

	block, _ := pem.Decode(certPEM)

	if block == nil || block.Type != "CERTIFICATE" {
		return "", fmt.Errorf("failed to decode PEM block containing certificate")
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

func SeederClient() (*mp.Client, error) {

	config := pkg.Config{
		Address:     YamlConfig.Router.Registry.Address,
		Fingerprint: YamlConfig.Router.Registry.Fingureprint,
		Timeout:     5 * time.Second,
	}

	client, err := mp.NewSdkOperation(config)
	if err != nil {
		return nil, fmt.Errorf("Seeder connection error %v", err)
	}
	return client, nil
}
