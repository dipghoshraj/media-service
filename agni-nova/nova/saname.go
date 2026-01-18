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

	defer conn.Close()

	sni, connbuffer, err := sni.PeekSNI(conn)
	if err != nil {
		log.Printf("Failed to extract SNI: %v", err)
		return
	}

	client, err := SeederClient()
	if err != nil {
		return
	}
	agent, err := client.GetAgentProxyMapping(context.Background(), "global", sni)
	if err != nil {
		log.Printf("Failed to get router for SNI %s: %v", sni, err)
		return
	}

	log.Println("[NOVA] found gateway", agent.GatewayAddress)

	backendConn, err := net.DialTimeout("tcp", agent.GatewayAddress, 3*time.Second)
	if err != nil {
		log.Printf("Failed to connect to backend %s: %v", agent.GatewayAddress, err)
		return
	}

	go io.Copy(backendConn, connbuffer)
	io.Copy(connbuffer, backendConn)

}
