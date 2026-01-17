package rpc

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/odio4u/agni-tunnels/agni-agent/pkg/bridge"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var rpcconn *grpc.ClientConn

func routerConnect(router string, gatewayIdentity string) *grpc.ClientConn {

	permfile := fmt.Sprintf("%s/client.pem", bridge.YamlConfig.Agent.Certs)
	permKeyfile := fmt.Sprintf("%s/client-key.pem", bridge.YamlConfig.Agent.Certs)
	clientCert, _ := tls.LoadX509KeyPair(permfile, permKeyfile)

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			cert, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return err
			}

			fp := sha256.Sum256(cert.Raw)
			expected := gatewayIdentity

			log.Println("current gateway connection data ", hex.EncodeToString(fp[:]))

			if hex.EncodeToString(fp[:]) != expected {
				return errors.New("server fingerprint mismatch")
			}
			return nil
		},
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient(
		router,
		grpc.WithTransportCredentials(creds),
	)

	if err != nil {
		panic(err)
	}

	rpcconn = conn

	return conn

}

func GetRouter() *grpc.ClientConn {
	return rpcconn
}
