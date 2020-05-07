package orderbook

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/types"
)

type (
	OrderBook interface {
		CreateOrder(context context.Context, userId uint64) (uint64, error)
		AddItem(context context.Context, orderId uint64) (uint64, error)
		RemoveItem(context context.Context, receiptItem uint64) error
		GeneratePaymentLink(context context.Context, paymentId uint64) error
		GetOrderByMessageId(ctx context.Context, messageId uint64) (order types.Order, err error)
		UpdateHintMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		UpdateMessageForOrder(ctx context.Context, orderId, messageId uint64) error
		UpdateOrderEditState(ctx context.Context, orderId uint64, state types.EditState) error
		GetActiveOrderForUser(ctx context.Context, userId uint64) (types.Order, error)
		GetActiveOrderMessageIdForUser(ctx context.Context, userId uint64) (uint64, error)
		UpdateReceiptItemName(ctx context.Context, name string, userId, receiptItemId uint64) (err error)
		UpdateReceiptItemPrice(ctx context.Context, price uint32, receiptItemId uint64) (err error)
		UpdateReceiptItemQty(ctx context.Context, qty int, receiptItemId uint64) (err error)
		UpdateCustomerEmail(ctx context.Context, email string, orderId uint64) (customerId uint64, err error)
		UpdateCustomerInstagram(ctx context.Context, instagram string, orderId uint64) (customerId uint64, err error)
		GetReceiptItem(ctx context.Context, receiptItemId uint64) (item types.ReceiptItem, err error)
		GetOrder(ctx context.Context, orderId uint64) (order types.Order, err error)
		SetActiveItemId(ctx context.Context, orderId uint64, receiptItemId uint64) error
		SetActivePaymentId(ctx context.Context, orderId uint64, paymentId uint64) error
		AddPayment(context context.Context, orderId uint64, method types.PaymentMethod) (uint64, error)
		RemovePayment(ctx context.Context, paymentId uint64) error
		RefundPayment(ctx context.Context, paymentId uint64, amount uint32) error
		UpdatePaymentAmount(ctx context.Context, paymentId uint64, amount uint32) error
		GetOrderMessage(ctx context.Context, messageId uint64) (types.OrderMessage, error)
		UpdateOrderMessageDisplayMode(ctx context.Context, messageId uint64, mode types.DisplayMode) error
		UpdateOrderState(ctx context.Context, id uint64, state types.OrderState) error
		GetPayment(ctx context.Context, paymentId uint64) (types.Payment, error)
		GetBankPayments(ctx context.Context) (payments []types.Payment, err error)
	}

	orderBook struct {
		storage storage.Storage
	}
)

func (ob orderBook) GetPayment(ctx context.Context, paymentId uint64) (types.Payment, error) {
	return ob.storage.GetPayment(ctx, paymentId)
}

func (ob orderBook) GetActiveOrderMessageIdForUser(ctx context.Context, userId uint64) (uint64, error) {
	return ob.storage.GetActiveOrderMessageIdForUser(ctx, userId)
}

//Returns new instance of bill maker
func NewOrderBook(storage storage.Storage) (OrderBook, error) {
	return orderBook{
		storage: storage,
	}, nil
}

//Returns new instance of bill maker
func MustNewOrderBook(storage storage.Storage) OrderBook {
	bm, err := NewOrderBook(storage)
	if err != nil {
		panic(err)
	}
	return bm
}

func (ob orderBook) UpdateOrderMessageDisplayMode(ctx context.Context, messageId uint64, mode types.DisplayMode) error {
	return ob.storage.UpdateOrderMessageDisplayMode(ctx, messageId, mode)
}

func (ob orderBook) GetOrderMessage(ctx context.Context, messageId uint64) (types.OrderMessage, error) {
	return ob.storage.GetOrderMessage(ctx, messageId)
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

func (ob orderBook) SetActivePaymentId(ctx context.Context, orderId uint64, paymentId uint64) error {
	return ob.storage.SetActivePaymentId(ctx, orderId, paymentId)
}

func (ob orderBook) UpdateCustomerEmail(ctx context.Context, email string, orderId uint64) (uint64, error) {
	return ob.storage.UpdateCustomerEmail(ctx, email, orderId)
}

func (ob orderBook) UpdateCustomerInstagram(ctx context.Context, instagram string, orderId uint64) (uint64, error) {
	return ob.storage.UpdateCustomerInstagram(ctx, instagram, orderId)
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

func (ob orderBook) RefundPayment(ctx context.Context, paymentId uint64, amount uint32) error {
	return ob.storage.RefundPayment(ctx, paymentId, amount)
}

func (ob orderBook) GetBankPayments(ctx context.Context) (payments []types.Payment, err error) {
	return ob.storage.GetBankPayments(ctx)
}
