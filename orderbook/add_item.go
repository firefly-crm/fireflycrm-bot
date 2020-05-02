package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/types"
)

//Adds receipt item to existing bill and returns it's unique id
func (b orderBook) AddItem(context context.Context, billId uint64, item types.Item) (uint64, error) {
	panic("implement me")
}
