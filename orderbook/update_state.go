package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/types"
)

func (ob orderBook) UpdateOrderEditState(ctx context.Context, orderId uint64, state types.EditState) error {
	return ob.storage.UpdateOrderEditState(ctx, orderId, state)
}

func (ob orderBook) UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error {
	return ob.storage.UpdateOrderState(ctx, orderId, state)
}
