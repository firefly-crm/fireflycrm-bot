package orderbook

import (
	"context"
)

//Adds receipt item to existing order and returns it's unique id
func (b orderBook) AddItem(ctx context.Context, orderId uint64) (uint64, error) {
	return b.storage.AddItemToOrder(ctx, orderId)
}
