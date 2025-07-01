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
		for _, vv := range v {
			req.writer.Header().Add(k, string(vv))
		}
	}
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
		return
	}

	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in StreamToClient: %v", r)
		}
		m.Close(id)
	}()

	for chunk := range req.bodyBuf {
		if req.writer == nil {
			log.Printf("Writer is nil for stream %s", id)
			return
		}

		_, err := req.writer.Write(chunk)
		if err != nil {
			log.Printf("Write error: %v", err)
			return
		}

		if f, ok := req.writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}
