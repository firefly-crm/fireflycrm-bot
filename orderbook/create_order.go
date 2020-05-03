package orderbook

import (
	"context"
)

func (b orderBook) CreateOrder(ctx context.Context, userId uint64) (uint64, error) {
	return b.storage.CreateOrder(ctx, userId)
}
