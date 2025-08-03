package connects

import (
	tunnelv1 "agni-cli/proto"
	"context"
	"fmt"
)

func (t *TunnelClient) readLoop(ctx context.Context,
	stream tunnelv1.TunnelService_TunnelStreamClient,
	out chan<- *tunnelv1.Envelope,
	errs chan<- error) {

	for {
		msg, err := stream.Recv()
		if err != nil {
			errs <- fmt.Errorf("recv error: %w", err)
			return
		}

		select {
		case out <- msg:
		case <-ctx.Done():
			return
		}
	}

}
