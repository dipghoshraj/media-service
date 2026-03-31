package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/config"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/rpc"
	"github.com/odio4u/agni-tunnels/agni-router/server"
	certpkg "github.com/odio4u/mem-sdk/certengine/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func gracefulShutdown(server *grpc.Server) {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Logger.Info("shutting down server")

	// Attempt graceful shutdown
	server.Stop()

}

// runGenCerts generates the router TLS certificates from values in agni-config.yaml
// and exits. Run with: agni-router gen-certs
func runGenCerts() {
	routerIps := []string{config.YamlConfig.Router.RouterIP}
	dns := []string{config.YamlConfig.Router.Dns}

	_, err := certpkg.GenerateSelfSignedGPR(config.YamlConfig.Router.Name, routerIps, dns)
	if err != nil {
		logger.Logger.Error("[Agni Router] failed to generate certificates", "error", err)
		os.Exit(1)
	}

	logger.Logger.Info("[Agni Router] certificates generated successfully",
		"name", config.YamlConfig.Router.Name,
		"ip", config.YamlConfig.Router.RouterIP,
		"dns", config.YamlConfig.Router.Dns,
	)
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "gen-certs" {
		runGenCerts()
		return
	}

	permfile := fmt.Sprintf("%s/server.pem", config.YamlConfig.Router.Certs)
	permfileKey := fmt.Sprintf("%s/server-key.pem", config.YamlConfig.Router.Certs)

	cert, err := tls.LoadX509KeyPair(permfile, permfileKey)
	if err != nil {
		logger.Logger.Error("failed to load server certificate", "error", err)
		os.Exit(1)
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

		VerifyPeerCertificate: config.AuthAgent,
	}

	fingurePrint, err := config.CertFingurePrint()
	if err != nil {
		logger.Logger.Error("failed to read certificate fingerprint", "error", err)
		os.Exit(1)
	}

	port := config.YamlConfig.Router.RouterPort

	port = fmt.Sprintf(":%s", port)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		logger.Logger.Error("failed to start gRPC listener", "listen_addr", port, "error", err)
		os.Exit(1)
	}

	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			stack := string(debug.Stack())
			logger.Logger.Error("panic recovered", "panic", fmt.Sprintf("%v", p), "stack", stack)
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
	go func() {
		server.RouterServer()
	}()

	logger.Logger.Info("gRPC server listening", "listen_addr", port)
	go func() {
		if err := s.Serve(lis); err != nil {
			logger.Logger.Error("gRPC server error", "error", err)
			os.Exit(1)
		}
	}()

	gracefulShutdown(s)

}
