package orderbook

import "context"

//Finalizes bill and sends it to bank. Returns payment link
func (b orderBook) GetPaymentLink(context context.Context, billId uint64) (string, error) {
	panic("implement me")
}
