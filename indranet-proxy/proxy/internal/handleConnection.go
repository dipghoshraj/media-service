package internal

import (
	"cosmo-proxy/internal/sni"
	"io"
	"log"
	"net"
	"time"
)

func HandleConnection(conn net.Conn) {
	// Extract the SNI from the connection

	defer conn.Close()

	// clientIP := GetClientIP(conn)

	sni, connbuffer, err := sni.PeekSNI(conn)
	if err != nil {
		log.Printf("Failed to extract SNI: %v", err)
		return
	}
	routerAddr, err := GetRouter(sni)
	if err != nil {
		log.Printf("Failed to get router for SNI %s: %v", sni, err)
		return
	}

	backendConn, err := net.DialTimeout("tcp", routerAddr, 3*time.Second)
	if err != nil {
		log.Printf("Failed to connect to backend %s: %v", routerAddr, err)
		return
	}

	go io.Copy(backendConn, connbuffer)
	io.Copy(connbuffer, backendConn)

}
