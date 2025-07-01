package inflight

import (
	pb "gateway-tunnel/proto"
	"log"
	"net/http"
)

func (m *inFlightManager) Register(id string, w http.ResponseWriter) {
	m.Lock()
	defer m.Unlock()

	if _, exists := m.requests[id]; exists {
		http.Error(w, "Request already in progress", http.StatusConflict)
		return
	}

	m.requests[id] = &inFlightRequest{
		writer:  w,
		headers: make(http.Header),
		bodyBuf: make(chan []byte, 1),
		done:    make(chan struct{}),
	}
}

func (m *inFlightManager) Resolve(id string, resp *pb.TunnelResponse) {

	m.Lock()

	req, ok := m.requests[id]
	if !ok {
		m.Unlock()
		log.Printf("Unknown response ID: %s", id)
		return
	}

	delete(m.requests, id)
	m.Unlock()

	for k, v := range resp.Headers {
		req.writer.Header().Set(k, v)
	}

	// for k, v := range resp.Headers {
	// 	for _, vv := range v {
	// 		req.writer.Header().Add(k, string(vv))
	// 	}
	// }
	req.writer.WriteHeader(int(resp.Status))
	_, _ = req.writer.Write(resp.Body)

	close(req.done)

}

func (m *inFlightManager) Close(id string) {
	m.Lock()
	req, ok := m.requests[id]
	if ok {
		delete(m.requests, id)
		close(req.bodyBuf)
		close(req.done)
	}
	m.Unlock()
}

func (m *inFlightManager) Stream(id string, chunk []byte) {
	m.RLock()
	req, ok := m.requests[id]
	m.RUnlock()
	if !ok {
		log.Printf("Unknown stream ID for streams: %s", id)
		return
	}
	select {
	case req.bodyBuf <- chunk:
	default:
		log.Printf("Stream buffer full for: %s", id)
	}
}

func (m *inFlightManager) GetDoneChan(id string) <-chan struct{} {
	m.RLock()
	defer m.RUnlock()
	return m.requests[id].done
}

func (m *inFlightManager) get(id string) (*inFlightRequest, bool) {
	m.RLock()
	defer m.RUnlock()

	req, ok := m.requests[id]
	if !ok {
		log.Printf("Unknown stream ID: %s", id)
		return nil, false
	}
	return req, true
}

func (m *inFlightManager) StreamToClient(id string) {
	req, ok := m.get(id)
	if !ok {
		log.Printf("Cannot stream: unknown stream ID: %s", id)
		return
	}
	defer m.Close(id)

	for chunk := range req.bodyBuf {
		if _, err := req.writer.Write(chunk); err != nil {
			log.Printf("Write to client failed for stream %s: %v", id, err)
			return
		}
		// Important to flush streamed data
		if f, ok := req.writer.(http.Flusher); ok {
			f.Flush()
		}
		log.Printf("Flushed %d bytes to client [stream: %s]", len(chunk), id)
	}
}
