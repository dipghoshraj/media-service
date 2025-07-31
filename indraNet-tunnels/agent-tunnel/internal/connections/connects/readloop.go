package connects

import (
	"agent-tunnel/proto"

	protobuf "google.golang.org/protobuf/proto"

	"context"
	"fmt"
	"log"
)

func (c *TunnelClient) readLoop(ctx context.Context, out chan<- *proto.Envelope, errs chan<- error) {
	// Placeholder for the read loop logic
	// This function should implement the logic to read messages from the connection
	// and handle them appropriately.

	// Example:
	// - Read messages from the connection
	// - Process each message based on its type
	// - Handle errors or close the connection if necessary

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			errs <- fmt.Errorf("read error: %w", err)
			return
		}

		var env proto.Envelope
		if err := protobuf.Unmarshal(data, &env); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		select {
		case out <- &env:
		case <-ctx.Done():
			return
		}
	}
}
