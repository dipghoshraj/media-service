package main

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
	"github.com/Purple-House/agni-tunnels/agni-router/pkg/config"
	"github.com/Purple-House/agni-tunnels/agni-router/pkg/rpc"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func gracefulShutdown(server *grpc.Server) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	// Attempt graceful shutdown
	server.Stop()

}

func main() {

	permfile := fmt.Sprintf("%s/server.pem", config.YamlConfig.Router.Certs)
	permfileKey := fmt.Sprintf("%s/server-key.pem", config.YamlConfig.Router.Certs)

	cert, err := tls.LoadX509KeyPair(permfile, permfileKey)
	if err != nil {
		log.Fatalf("failed to load server certificate: %v", err)
	}

	servertLs := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
		MaxVersion:   tls.VersionTLS13,

		CipherSuites: []uint16{
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
		},

		SessionTicketsDisabled:   true,
		PreferServerCipherSuites: true,

		CurvePreferences: []tls.CurveID{
			tls.X25519,
			tls.CurveP256,
		},
		Renegotiation:      tls.RenegotiateNever,
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAnyClientCert,

		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
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

			if hex.EncodeToString(fp[:]) != "565d3b03ef609573f47329e7896a9ff825f3bb37eb92f0f2383029e895adfca5" {
				return errors.New("client fingerprint mismatch")
			}
			return nil
		},
	}

	fingurePrint, err := config.CertFingurePrint()
	if err != nil {
		log.Fatalf("Failed to print certificate fingerprint: %v", err)
	}

	port := config.YamlConfig.Router.RouterPort

	port = fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			stack := string(debug.Stack())
			log.Printf("[PANIC RECOVERED] %v\nSTACK TRACE:\n%s", p, stack)
			return fmt.Errorf("internal server error")
		}),
	}

	s := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(servertLs)),
		grpc.UnaryInterceptor(grpc_recovery.UnaryServerInterceptor(recoveryOpts...)),
		grpc.StreamInterceptor(grpc_recovery.StreamServerInterceptor(recoveryOpts...)),
	)

	err = config.SeedGatway(fingurePrint)

	tunnel.RegisterAgniTunnelServer(s, &rpc.TunnelRpc{})

	// Start the server
	log.Println("Server is running on port ", port)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	gracefulShutdown(s)

}
