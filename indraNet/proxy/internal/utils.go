package internal

import (
	"net"
	"strings"
)

func GetClientIP(conn net.Conn) string {
	addr := conn.RemoteAddr().String()
	parts := strings.Split(addr, ":")
	return parts[0]
}

// func GetLimiter(ip string) *rate.Limiter {
// 	if limiterIface, ok := ipLimiters.Load(ip); ok {
// 		return limiterIface.(*rate.Limiter)
// 	}
// 	limiter := rate.NewLimiter(10, 20) // 10 req/sec, burst 20
// 	ipLimiters.Store(ip, limiter)
// 	return limiter
// }
