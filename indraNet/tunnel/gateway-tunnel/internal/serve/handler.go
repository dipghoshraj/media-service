package serve

import (
	"bytes"
	"fmt"
	"gateway-tunnel/internal/client/inflight"
	"gateway-tunnel/internal/client/session"
	"io"
	"log"
	"net/http"
	"time"

	pb "gateway-tunnel/proto"

	"github.com/google/uuid"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	appID := extractDomain(r.Host)

	agent, ok := session.Registry.GetSession(appID)
	if !ok {
		http.Error(w, "Agent session not found", http.StatusNotFound)
		return
	}

	var body bytes.Buffer

	fmt.Fprintf(&body, "%s %s HTTP/1.1\r\n", r.Method, r.URL.RequestURI())
	if r.Host != "" {
		fmt.Fprintf(&body, "Host: %s\r\n", r.Host)
	} else {
		fmt.Fprintf(&body, "Host: localhost\r\n")
	}

	// End of headers
	body.WriteString("\r\n")

	// If method has body, copy it
	if r.Body != nil && (r.Method == "POST" || r.Method == "PUT") {
		io.Copy(&body, r.Body)
	}

	defer r.Body.Close()

	id := uuid.New().String()
	tunnelReq := &pb.TunnelRequest{
		Id:      id,
		Method:  r.Method,
		Path:    r.URL.RequestURI(),
		Body:    body.Bytes(),
		Headers: flattenHeaders(r.Header),
	}

	inflight.InFlightManager.Register(id, w)
	go inflight.InFlightManager.StreamToClient(id)

	select {
	case agent.SendChan <- tunnelReq:
		log.Printf("Sent request %s to agent %s", id, appID)
	default:
		http.Error(w, "Agent busy", http.StatusTooManyRequests)
		inflight.InFlightManager.Close(id)
		return
	}

	select {
	case <-time.After(10 * time.Second):
		inflight.InFlightManager.Close(id)
		http.Error(w, "Timeout", http.StatusGatewayTimeout)
	case <-inflight.InFlightManager.GetDoneChan(id):
		// Response already written
	}

}
