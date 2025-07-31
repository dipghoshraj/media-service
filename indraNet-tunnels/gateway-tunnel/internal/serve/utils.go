package serve

import (
	"net/http"
	"strings"
)

func extractDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

func flattenHeaders(hdr http.Header) map[string]string {
	out := make(map[string]string)
	for k, vals := range hdr {
		if len(vals) > 0 {
			out[k] = vals[0]
		}
	}
	return out
}
