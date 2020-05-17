package orderbook

import (
	"context"
	"database/sql"
)

func (ob orderBook) UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error {
	msgId := sql.NullInt64{Valid: messageId != 0, Int64: int64(messageId)}
	return ob.storage.UpdateHintMessageForOrder(ctx, orderId, msgId)
}
