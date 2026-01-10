package rpc

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var rpcconn *grpc.ClientConn

func routerConnect(router string, gatewayIdentity string) *grpc.ClientConn {

	tlsConfig := &tls.Config{
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
