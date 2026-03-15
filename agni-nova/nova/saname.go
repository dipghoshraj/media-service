package nova

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/odio4u/mem-sdk/sni"
)

// func saname() {
// 	serverName, _, err := sni.PeekSNI(conn)
// }

func HandleStream(conn net.Conn) {

	start := time.Now()
	clientAddr := conn.RemoteAddr().String()
	log.Printf("[NOVA] connection from %s", clientAddr)

	defer conn.Close()

	serverName, connbuffer, err := sni.PeekSNI(conn)
	if err != nil {
		if err != io.EOF {
			log.Printf("[NOVA] failed to extract SNI from %s: %v", clientAddr, err)
			return
		}
	}
	log.Printf("[NOVA] SNI=%s client=%s", serverName, clientAddr)

	client, err := SeederClient()
	if err != nil {
		return
	}
	agent, err := client.GetAgentProxyMapping(context.Background(), "global", serverName)
	if err != nil {
		log.Printf("[NOVA] failed to get router for SNI=%s client=%s: %v", serverName, clientAddr, err)
		return
	}

	log.Printf("[NOVA] routing SNI=%s client=%s → gateway=%s", serverName, clientAddr, agent.GatewayAddress)

	backendConn, err := net.DialTimeout("tcp", agent.GatewayAddress, 3*time.Second)
	if err != nil {
		log.Printf("[NOVA] failed to connect to backend %s: %v", agent.GatewayAddress, err)
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

	duration := time.Since(start)
	if toClientErr != nil || toBackendErr != nil {
		log.Printf("[NOVA] relay error SNI=%s client=%s duration=%v err(backend→client)=%v err(client→backend)=%v",
			serverName, clientAddr, duration.Round(time.Millisecond), toClientErr, toBackendErr)
		if duration < 2*time.Second {
			log.Printf("[NOVA] WARNING: relay for %s closed in %v with an error — client likely rejected the server TLS certificate. Use curl.exe -k to skip cert validation.",
				clientAddr, duration.Round(time.Millisecond))
		}
	} else {
		log.Printf("[NOVA] relay closed SNI=%s client=%s duration=%v", serverName, clientAddr, duration.Round(time.Millisecond))
	}

}
