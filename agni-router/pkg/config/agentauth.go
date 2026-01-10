package config

// import (
// 	"crypto/sha256"
// 	"crypto/x509"
// 	"encoding/hex"
// 	"errors"
// 	"log"
// )

// func (rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
// 	if len(rawCerts) == 0 {
// 		return errors.New("no client certificate provided")
// 	}

// 	clientCert, err := x509.ParseCertificate(rawCerts[0])
// 	if err != nil {
// 		return err
// 	}

// 	fp := sha256.Sum256(clientCert.Raw)
// 	log.Println("Client fingerprint:", hex.EncodeToString(fp[:]))

// 	if hex.EncodeToString(fp[:]) != "expectedClientFingerprint" {
// 		return errors.New("client fingerprint mismatch")
// 	}
// 	return nil
// }
