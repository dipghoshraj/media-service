package rpc

import (
	"log"

	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/session"
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

	domain, exist := session.Seeder.GetDomainMap(req.AgentId)
	if !exist {
		log.Panic("[Agni Router] can not found id domain mapping", req.AgentId)
	}

	agentSession := &session.AgentSession{
		AppID:  req.AgentId,
		Stream: &stream,
	}
	session.Seeder.Register(domain, agentSession)

	if err = stream.Send(ackMessage); err != nil {
		return status.Error(codes.ResourceExhausted, "Agent ackoledgement failed")
	}

	log.Printf("Agent %s connected successfully", req.AgentId)

	select {} // temporary
}

func checkConnect(req *tunnel.ConnectRequest) bool {
	log.Println("check the agent", req.AgentId)
	// log.Println("check the token", req.Token)
	// log.Println("check the signature", req.Signature)
	// log.Println("check the none", req.Nonce)

	return true
}
