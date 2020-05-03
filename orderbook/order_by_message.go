package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/types"
)

func (o orderBook) GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error) {
	return o.storage.GetOrderByMessageId(ctx, messageId)
}
