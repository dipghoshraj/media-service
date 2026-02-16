package message

import (
	"context"

	"github.com/google/uuid"
)

func (m *Message) ConnectAck(ctx context.Context) error {

	connection_id := uuid.New().String()

	connctx := NewConnectionCtx(connection_id)
	ctx = context.WithValue(ctx, "conn", connctx)
	return nil
}
