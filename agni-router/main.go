package main

import (
	"crypto/tls"
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
		Renegotiation: tls.RenegotiateNever,
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
