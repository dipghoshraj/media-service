package serve

import (
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	id := uuid.New().String()
	tunnelReq := &pb.TunnelRequest{
		Id:      id,
		Method:  r.Method,
		Path:    r.URL.RequestURI(),
		Body:    body,
		Headers: flattenHeaders(r.Header),
	}

	inflight.InFlightManager.Register(id, w)

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
