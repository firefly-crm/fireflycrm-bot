package orderbook

import "context"

func (ob orderBook) UpdateMessageForOrder(ctx context.Context, orderId, messageId uint64) error {
	return ob.storage.UpdateMessageForOrder(ctx, orderId, messageId)
}