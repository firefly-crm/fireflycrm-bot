package orderbook

import "context"

func (ob orderBook) UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error {
	return ob.storage.UpdateHintMessageForOrder(ctx, orderId, messageId)
}
