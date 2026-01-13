package sni

import (
	"encoding/binary"
	"errors"
)

func parseTLSRecord(data []byte) ([]byte, error) {
	if len(data) < 5 || data[0] != 0x16 {
		return nil, errors.New("not a TLS handshake record")
	}

	recordLen := int(binary.BigEndian.Uint16(data[3:5]))
	if len(data) < 5+recordLen {
		return nil, errors.New("incomplete TLS record")
	}

	return data[5 : 5+recordLen], nil
}

func parseSNI(data []byte) (string, error) {
	if len(data) < 2 {
		return "", errors.New("invalid SNI extension")
	}

	listLen := int(binary.BigEndian.Uint16(data))
	data = data[2:]

	if len(data) < listLen {
		return "", errors.New("SNI list truncated")
	}

	for len(data) >= 3 {
		nameType := data[0]
		nameLen := int(binary.BigEndian.Uint16(data[1:3]))
		data = data[3:]

		if len(data) < nameLen {
			return "", errors.New("invalid SNI name length")
		}

		if nameType == 0x00 {
			return string(data[:nameLen]), nil
		}

		data = data[nameLen:]
	}

	return "", errors.New("no DNS SNI found")
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

func SniStream(data []byte) (string, error) {

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
