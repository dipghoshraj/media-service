package rpc

import (
	tunnel "github.com/odio4u/agni-schema/tunnel"
	"github.com/odio4u/agni-tunnels/agni-router/pkg/logger"
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
		logger.Logger.Error("agent id to domain mapping not found", "agent_id", req.AgentId)
		return status.Errorf(codes.NotFound, "domain mapping not found for agent %s", req.AgentId)
	}

	agentSession := &session.AgentSession{
		AppID:  req.AgentId,
		Stream: &stream,
	}
	session.Seeder.Register(domain, agentSession)

	if err = stream.Send(ackMessage); err != nil {
		return status.Error(codes.ResourceExhausted, "Agent ackoledgement failed")
	}

	logger.Logger.Info("agent connected", "agent_id", req.AgentId, "domain", domain)

	select {} // temporary
}

func checkConnect(req *tunnel.ConnectRequest) bool {
	logger.Logger.Info("validating connect request", "agent_id", req.AgentId)
	return true
}
