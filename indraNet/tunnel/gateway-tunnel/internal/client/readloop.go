package client

import (
	"gateway-tunnel/internal/client/session"
	pb "gateway-tunnel/proto"
	"log"

	"gateway-tunnel/internal/client/inflight"

	"google.golang.org/protobuf/proto"
)

func readLoop(agenetSession *session.AgentSession, done chan struct{}) {

	defer func() {
		log.Println("Read loop ended for session:", agenetSession.AppID)
		session.Registry.Unregister(agenetSession.AppID)
		agenetSession.Conn.Close()
	}()

	for {
		_, msgByte, err := agenetSession.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read error for session %s: %v", agenetSession.AppID, err)
			done <- struct{}{}
			return
		}

		var env pb.Envelope
		if err := proto.Unmarshal(msgByte, &env); err != nil {
			log.Printf("Proto unmarshal error: %v", err)
			continue
		}

		switch msg := env.Message.(type) {
		case *pb.Envelope_Response:
			inflight.InFlightManager.Resolve(msg.Response.Id, msg.Response)
		case *pb.Envelope_Data:
			inflight.InFlightManager.Stream(msg.Data.Id, msg.Data.Chunk)
		case *pb.Envelope_Close:
			inflight.InFlightManager.Close(msg.Close.Id)
		// case *pb.Envelope_Control:
		//     handleControl(msg.Control)
		default:
			log.Printf("Unknown message type")
		}
	}
}
