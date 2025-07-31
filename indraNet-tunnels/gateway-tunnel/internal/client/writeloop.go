package client

import (
	"gateway-tunnel/internal/client/session"
	"log"

	pb "gateway-tunnel/proto"

	"github.com/gorilla/websocket"
	protobuf "google.golang.org/protobuf/proto"
)

func writeLoop(sess *session.AgentSession, done chan struct{}) {
	for req := range sess.SendChan {
		env := &pb.Envelope{
			Message: &pb.Envelope_Request{Request: req},
		}

		data, err := protobuf.Marshal(env)
		if err != nil {
			log.Printf("Marshal error: %v", err)
			continue
		}

		if err := sess.Conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			log.Printf("Agent [%s] write error: %v", sess.AppID, err)
			done <- struct{}{}
			return
		}
	}
}
