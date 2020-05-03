package orderbook

import "context"

//Removes item from created order
func (b orderBook) RemoveItem(ctx context.Context, itemId uint64) error {
	return b.storage.RemoveReceiptItem(ctx, itemId)
}
