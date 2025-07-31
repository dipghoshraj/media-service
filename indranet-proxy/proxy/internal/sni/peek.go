package sni

import (
	"net"
	"strings"
)

func PeekSNI(conn net.Conn) (string, net.Conn, error) {
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		return "", nil, err
	}
	serverName, err := ExtractHostFromHTTPRequest(buf[:n])
	// serverName, err := ExtractSNIFromRequest(buf[:n]) this is for TLS SNI extraction in current context request is flowing over
	domain := strings.Split(serverName, ".")
	if len(domain) == 0 {
		return "", nil, nil // No SNI found
	}

	if err != nil {
		return "", nil, err
	}

	return domain[0], &ConnBuffer{buf: buf[:n], conn: conn}, nil
}
