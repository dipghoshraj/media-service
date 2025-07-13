package sni

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

func ExtractHostFromHTTPRequest(data []byte) (string, error) {
	lines := strings.Split(string(data), "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Host:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Host:")), nil
		}
	}
	return "", fmt.Errorf("host header not found in request")
}

func ExtractSNIFromRequest(data []byte) (string, error) {

	record, err := parseTLSRecord(data)
	if err != nil {
		return "", err
	}

	hello, err := parseHandshake(record)
	if err != nil {
		return "", err
	}

	hello, err = parseSessionID(hello)
	if err != nil {
		return "", err
	}

	hello, err = parseCipherSuites(hello)
	if err != nil {
		return "", err
	}

	hello, err = parseCompressionMethods(hello)
	if err != nil {
		return "", err
	}

	extensions, err := parseExtensions(hello)
	if err != nil {
		return "", err
	}

	return findSNIExtension(extensions)

}

func parseTLSRecord(data []byte) ([]byte, error) {
	if len(data) < 5 || data[0] != 0x16 {
		return nil, errors.New("not a TLS handshake record")
	}
	return data[5:], nil
}

func parseHandshake(data []byte) ([]byte, error) {
	if len(data) < 4 || data[0] != 0x01 {
		return nil, errors.New("not a ClientHello")
	}

	data = data[4:] // skip handshake header
	if len(data) < 34 {
		return nil, errors.New("incomplete ClientHello header")
	}

	return data[34:], nil // skip version + random
}

func parseSessionID(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, errors.New("missing session ID length")
	}

	sidLen := int(data[0])
	data = data[1:]

	if len(data) < sidLen {
		return nil, errors.New("truncated session ID")
	}

	return data[sidLen:], nil
}

func parseCipherSuites(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return nil, errors.New("missing cipher suites length")
	}

	cipherLen := int(binary.BigEndian.Uint16(data))
	data = data[2:]

	if len(data) < cipherLen {
		return nil, errors.New("truncated cipher suites")
	}

	return data[cipherLen:], nil
}

func parseCompressionMethods(data []byte) ([]byte, error) {
	if len(data) < 1 {
		return nil, errors.New("missing compression methods length")
	}

	compLen := int(data[0])
	data = data[1:]

	if len(data) < compLen {
		return nil, errors.New("truncated compression methods")
	}

	return data[compLen:], nil
}

func parseExtensions(data []byte) ([]byte, error) {
	if len(data) < 2 {
		return nil, errors.New("missing extensions length")
	}

	extLen := int(binary.BigEndian.Uint16(data))
	data = data[2:]

	if len(data) < extLen {
		return nil, errors.New("truncated extensions")
	}

	return data[:extLen], nil
}

func findSNIExtension(data []byte) (string, error) {
	for len(data) >= 4 {
		extType := binary.BigEndian.Uint16(data)
		extLen := int(binary.BigEndian.Uint16(data[2:]))
		data = data[4:]

		if len(data) < extLen {
			return "", errors.New("truncated extension")
		}

		if extType == 0x00 {
			return parseSNI(data[:extLen])
		}

		data = data[extLen:]
	}

	return "", errors.New("SNI not found")
}

func parseSNI(data []byte) (string, error) {
	if len(data) < 5 {
		return "", errors.New("invalid SNI extension length")
	}

	listLen := int(binary.BigEndian.Uint16(data))
	if len(data) < 2+listLen {
		return "", errors.New("SNI list truncated")
	}

	if data[2] != 0x00 {
		return "", errors.New("unsupported SNI name type")
	}

	nameLen := int(binary.BigEndian.Uint16(data[3:5]))
	if 5+nameLen > len(data) {
		return "", errors.New("invalid SNI name length")
	}

	return string(data[5 : 5+nameLen]), nil
}
