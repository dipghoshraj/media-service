package inflight

import (
	"net/http"
	"sync"
)

type inFlightRequest struct {
	writer  http.ResponseWriter
	headers http.Header
	bodyBuf chan []byte
	done    chan struct{}
}

type inFlightManager struct {
	sync.RWMutex
	requests map[string]*inFlightRequest
}

var InFlightManager = &inFlightManager{
	requests: make(map[string]*inFlightRequest),
}
