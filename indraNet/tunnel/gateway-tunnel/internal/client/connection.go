package client

import (
	"log"
	"net/http"
	"time"

	pb "gateway-tunnel/proto"

	"gateway-tunnel/internal/client/session"

	"github.com/gorilla/websocket"
	protobuf "google.golang.org/protobuf/proto"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println("Handshake failed:", err)
		return
	}

	var envalop pb.Envelope
	if err := protobuf.Unmarshal(msgBytes, &envalop); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid envelope"))
		conn.Close()
		log.Println("Failed to unmarshal envelope:", err)
		return
	}

	connectMsg, ok := envalop.Message.(*pb.Envelope_Connect)
	if !ok {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid connect message"))
		conn.Close()
		log.Println("Invalid connect message type")
		return
	}

	req := connectMsg.Connect
	if !verifyToken(req.Token, req.Signature) {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid token or signature"))
		conn.Close()
		log.Println("Invalid token or signature")
		return
	}

	agentSession := &session.AgentSession{
		AppID:    req.AgentId,
		Conn:     conn,
		SendChan: make(chan *pb.TunnelRequest, 10),
		LastSeen: time.Now(),
	}

	session.Registry.Register(req.AgentId, agentSession)

	log.Printf("Agent [%s] connected", req.AgentId)
	// go readLoop(agentSession)
	// go writeLoop(agentSession)

	select {}
}
