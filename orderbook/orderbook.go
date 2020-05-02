package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/DarthRamone/modulbank-go"
)

type (
	OrderBook interface {
		//Creates bill and returns its identifier
		CreateOrder(context context.Context, opts types.OrderOptions) (uint64, error)
		//Adds receipt item to existing bill and returns it's unique id
		AddItem(context context.Context, billId uint64, item types.Item) (uint64, error)
		//Removes item from created bill
		RemoveItem(context context.Context, itemId uint64) error
		//Finalizes bill and sends it to bank. Returns payment link
		GetPaymentLink(context context.Context, billId uint64) (string, error)
	}

	orderBook struct {
		storage   storage.Storage
		modulBank modulbank.API
	}
)

//Returns new instance of bill maker
func NewOrderBook(storage storage.Storage, modulBank modulbank.API) (OrderBook, error) {
	return orderBook{
		storage:   storage,
		modulBank: modulBank,
	}, nil
}

//Returns new instance of bill maker
func MustNewOrderBook(storage storage.Storage, modulBank modulbank.API) OrderBook {
	bm, err := NewOrderBook(storage, modulBank)
	if err != nil {
		panic(err)
	}
	return bm
}
