package rpc

import (
	"log"

	tunnel "github.com/Purple-House/agni-schema/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *TunnelRpc) Connect(stream tunnel.AgniTunnel_ConnectServer) error {

	envalop, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal, "failed to receive connect request: %v", err)

	}

	req := envalop.GetConnect()
	if req == nil {
		return status.Error(codes.InvalidArgument, "first message must be ConnectRequest")
	}

	validconnect := checkConnect(req)
	if !validconnect {
		return status.Error(codes.Aborted, "fuck off")
	}

	ackMessage := &tunnel.Envelope{
		Message: &tunnel.Envelope_ConnectAck{
			ConnectAck: &tunnel.ConnectAck{
				AgentId:  req.AgentId,
				Accepted: true,
			},
		},
	}

	if err = stream.Send(ackMessage); err != nil {
		return status.Error(codes.ResourceExhausted, "Agent ackoledgement failed")
	}

	log.Printf("Agent %s connected successfully", req.AgentId)

	return nil
}

func checkConnect(req *tunnel.ConnectRequest) bool {
	log.Println("check the agent", req.AgentId)
	log.Println("check the token", req.Token)
	log.Println("check the signature", req.Signature)
	log.Println("check the none", req.Nonce)

	return true
}
