package nova

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/odio4u/mem-sdk/sni"
)

func HandleStream(conn net.Conn) {
	start := time.Now()
	clientAddr := conn.RemoteAddr().String()
	Logger.Info("connection accepted", "client_addr", clientAddr)

	defer conn.Close()

	serverName, connbuffer, err := sni.PeekSNI(conn)
	if err != nil {
		if err != io.EOF {
			Logger.Error("failed to extract SNI", "client_addr", clientAddr, "error", err)
			return
		}
	}
	Logger.Info("SNI extracted", "sni", serverName, "client_addr", clientAddr)

	client, err := SeederClient()
	if err != nil {
		Logger.Error("seeder client error", "error", err)
		return
	}
	agent, err := client.GetAgentProxyMapping(context.Background(), "global", serverName)
	if err != nil {
		Logger.Error("failed to get router", "sni", serverName, "client_addr", clientAddr, "error", err)
		return
	}

	Logger.Info("routing connection", "sni", serverName, "client_addr", clientAddr, "gateway", agent.GatewayAddress)

	backendConn, err := net.DialTimeout("tcp", agent.GatewayAddress, 3*time.Second)
	if err != nil {
		Logger.Error("failed to connect to backend", "gateway", agent.GatewayAddress, "error", err)
		return
	}
	defer backendConn.Close()

	errc := make(chan error, 1)
	go func() {
		_, err := io.Copy(backendConn, connbuffer)
		errc <- err
	}()
	_, toClientErr := io.Copy(connbuffer, backendConn)
	toBackendErr := <-errc

	durationMs := duration(start)
	if toClientErr != nil || toBackendErr != nil {
		Logger.Error("relay closed with errors",
			"sni", serverName,
			"client_addr", clientAddr,
			"duration_ms", durationMs,
			"err_backend_to_client", toClientErr,
			"err_client_to_backend", toBackendErr,
		)
		if durationMs < 2000 {
			Logger.Warn("relay closed quickly — client may have rejected the server TLS certificate",
				"client_addr", clientAddr,
				"duration_ms", durationMs,
				"hint", "use -k / --insecure flag to skip cert validation during testing",
			)
		}
	} else {
		Logger.Info("relay closed",
			"sni", serverName,
			"client_addr", clientAddr,
			"duration_ms", durationMs,
		)
	}
}

// duration returns elapsed milliseconds since start as an int64.
func duration(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
