package config

import (
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mp "github.com/odio4u/mem-sdk/memsdk/maps"
	"github.com/odio4u/mem-sdk/memsdk/pkg"
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
		Address:     YamlConfig.Router.Seeder.Address,
		Fingerprint: YamlConfig.Router.Seeder.Fingureprint,
		Timeout:     5 * time.Second,
	}

	client, err := mp.NewSdkOperation(config)
	if err != nil {
		return nil, fmt.Errorf("Seeder connection error %v", err)
	}
	return client, nil
}

func SeedGatway(fingureprint string) error {
	client, err := SeederClient()
	if err != nil {
		return err
	}

	defer client.Close()

	region := "global"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	port, _ := strconv.ParseInt(YamlConfig.Router.RouterPort, 10, 32)
	proxy_port, _ := strconv.ParseInt(YamlConfig.Router.ProxtPort, 10, 32)

	router := mp.AddRouterRequest{
		Region:     region,
		RouterIp:   YamlConfig.Router.RouterIP,
		RouterPort: int32(proxy_port),
		Identity:   fingureprint,
		RpcPort:    int32(port),
	}

	gateways, err := client.Addgateway(ctx, router)
	if err != nil {
		panic(err)
	}
	log.Println("Added Gateway:", gateways.ID, gateways.GatewayPort, gateways.WssPort, gateways.IP)
	return nil
}
