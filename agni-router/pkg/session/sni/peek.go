package sni

import (
	"log"
	"net"
)

func PeekSNI(conn net.Conn) (string, net.Conn, error) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", nil, err
	}
	serverName, err := SniStream(buf[:n])

	log.Printf("found the server name domain: %s", serverName)

	if err != nil {
		return "", nil, err
	}

	return serverName, &ConnBuffer{buf: buf[:n], conn: conn}, nil
}
