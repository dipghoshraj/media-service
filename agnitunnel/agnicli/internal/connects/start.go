package connects

import (
	"context"
	"fmt"
	"log"
	"time"

	tunnelv1 "agni-cli/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (t *TunnelClient) Start(ctx context.Context) error {
	for {
		err := t.runSession(ctx)
		if err != nil {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second): // TODO: use exponential backoff
			}
		} else {
			return nil // Exit if runSession completes successfully
		}
	}
}

func (t *TunnelClient) runSession(ctx context.Context) error {
	// this function in not only for running the grpc connecto
	// in near future we are going to use this for the discovery of the
	// proximity gateway to keep the agent more dev friendly and easy to use

	conn, err := grpc.NewClient(t.Cfg.GatewayURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to create %s gRPC client: %v", t.Cfg.GatewayURL, err)
		return fmt.Errorf("failed to create gRPC client: %w", err)
	}

	defer conn.Close()
	defer func() {
		if t.stream != nil {
			t.stream.CloseSend()
			log.Println("Stream closed.")
		}
	}()

	client := tunnelv1.NewTunnelServiceClient(conn)
	stream, err := client.TunnelStream(ctx)
	if err != nil {
		log.Printf("Failed to start TunnelStream: %v", err)
		return fmt.Errorf("failed to start TunnelStream: %w", err)
	}
	t.stream = stream

	req := &tunnelv1.Envelope{
		Message: &tunnelv1.Envelope_ConnectRequest{
			ConnectRequest: &tunnelv1.ConnectRequest{
				AgentId:   t.Cfg.AgentID,
				Token:     t.Cfg.Token,
				Timestamp: time.Now().Unix(),
				Nonce:     GenerateNonce(16), // Add real nonce
				Signature: GenerateNonce(32), // Add real signature
			},
		},
	}
	if err := t.stream.Send(req); err != nil {
		return fmt.Errorf("failed to send connect request: %w", err)
	}

	log.Println("Connected to gateway successfully.")

	msgs := make(chan *tunnelv1.Envelope, 10)
	errs := make(chan error, 2)

	go t.readLoop(ctx, t.stream, msgs, errs)

	for {
		select {
		case <-msgs:
			//TODO: Handle incoming messages
		case err := <-errs:
			if err != nil {
				log.Printf("Error in stream: %v", err)
				return fmt.Errorf("stream error: %w", err)
			}
		case <-ctx.Done():
			log.Printf("Context canceled or timeout: %v", ctx.Err())
			return ctx.Err()
		}
	}

}
