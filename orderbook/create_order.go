package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/types"
)

//Creates bill and returns its identifier
func (b orderBook) CreateOrder(ctx context.Context, opts types.OrderOptions) (uint64, error) {
	return b.storage.CreateOrder(ctx, opts)
}
