package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/DarthRamone/modulbank-go"
)

type (
	OrderBook interface {
		CreateOrder(context context.Context, userId uint64) (uint64, error)
		AddItem(context context.Context, orderId uint64) (uint64, error)
		RemoveItem(context context.Context, receiptItem uint64) error
		GetPaymentLink(context context.Context, billId uint64) (string, error)
		GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error)
		UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		UpdateMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		UpdateOrderState(ctx context.Context, orderId uint64, state types.OrderState) error
		GetActiveOrderForUser(ctx context.Context, userId uint64) (types.Order, error)
		UpdateReceiptItemName(ctx context.Context, name string, userId, receiptItemId uint64) (err error)
		UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error)
		UpdateReceiptItemQty(ctx context.Context, qty int, receiptItemId uint64) (err error)
		UpdateCustomerEmail(ctx context.Context, email string, orderId uint64) (customerId uint64, err error)
		GetReceiptItem(ctx context.Context, receiptItemId uint64) (item types.ReceiptItem, err error)
		GetOrder(ctx context.Context, orderId uint64) (order types.Order, err error)
		SetActiveItemId(ctx context.Context, orderId uint64, receiptItemId uint64) error
		AddPayment(context context.Context, orderId uint64, method types.PaymentMethod) (uint64, error)
		RemovePayment(ctx context.Context, paymentId uint64) error
		UpdatePaymentAmount(ctx context.Context, paymentId uint64, amount uint32) error
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

func (ob orderBook) UpdateReceiptItemName(ctx context.Context, name string, userId, receiptItemId uint64) (err error) {
	return ob.storage.UpdateReceiptItemName(ctx, name, userId, receiptItemId)
}

func (ob orderBook) UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error) {
	return ob.storage.UpdateReceiptItemPrice(ctx, price, receiptItemId)
}

func (ob orderBook) UpdateReceiptItemQty(ctx context.Context, qty int, receiptItemId uint64) (err error) {
	return ob.storage.UpdateReceiptItemQty(ctx, qty, receiptItemId)
}

func (ob orderBook) GetReceiptItem(ctx context.Context, receiptItemId uint64) (item types.ReceiptItem, err error) {
	return ob.storage.GetReceiptItem(ctx, receiptItemId)
}

func (ob orderBook) GetOrder(ctx context.Context, orderId uint64) (order types.Order, err error) {
	return ob.storage.GetOrder(ctx, orderId)
}

func (ob orderBook) SetActiveItemId(ctx context.Context, orderId uint64, receiptItemId uint64) error {
	return ob.storage.SetActiveItemId(ctx, orderId, receiptItemId)
}

func (ob orderBook) UpdateCustomerEmail(ctx context.Context, email string, orderId uint64) (uint64, error) {
	return ob.storage.UpdateCustomerEmail(ctx, email, orderId)
}

func (ob orderBook) AddPayment(context context.Context, orderId uint64, method types.PaymentMethod) (uint64, error) {
	return ob.storage.AddPayment(context, orderId, method)
}

func (ob orderBook) RemovePayment(ctx context.Context, paymentId uint64) error {
	return ob.storage.RemovePayment(ctx, paymentId)
}

func (ob orderBook) UpdatePaymentAmount(ctx context.Context, paymentId uint64, amount uint32) error {
	return ob.storage.UpdatePaymentAmount(ctx, paymentId, amount)
}
