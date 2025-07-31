package client

import "log"

func verifyToken(token, secret string) bool {
	log.Printf("Verifying token: %s with secret: %s", token, secret)
	return true
}
