package client

import (
	"gateway-tunnel/internal/client/session"
	pb "gateway-tunnel/proto"
	"log"

	"google.golang.org/protobuf/proto"
)

func readLoop(agenetSession *session.AgentSession) {

	defer func() {
		log.Println("Read loop ended for session:", agenetSession.AppID)
		session.Registry.Unregister(agenetSession.AppID)
		agenetSession.Conn.Close()
	}()

	for {
		_, msgByte, err := agenetSession.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for session %s: %v", agenetSession.AppID, err)
			return
		}

		var env pb.Envelope
		if err := proto.Unmarshal(msgByte, &env); err != nil {
			log.Printf("Proto unmarshal error: %v", err)
			continue
		}

		switch msg := env.Message.(type) {
		case *pb.Envelope_Response:
			log.Printf("Received response for session %s: %s", agenetSession.AppID, msg.Response.Id)
		default:
			log.Printf("Received unknown message type for session %s: %T", agenetSession.AppID, msg)
		}
	}
}
