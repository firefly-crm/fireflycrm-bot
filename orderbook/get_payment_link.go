package orderbook

import (
	"context"
	"fmt"
	mb "github.com/DarthRamone/modulbank-go"
	"github.com/brianvoe/gofakeit"
	"net/http"
	"strconv"
	"time"
)

func (b orderBook) GeneratePaymentLink(ctx context.Context, paymentId uint64) error {
	payment, err := b.storage.GetPayment(ctx, paymentId)
	if err != nil {
		return fmt.Errorf("failed to get payment: %w", err)
	}

	order, err := b.storage.GetOrder(ctx, payment.OrderId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	user, err := b.storage.GetUser(ctx, order.UserId)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	customer, err := b.storage.GetCustomer(ctx, uint64(order.CustomerId.Int64))
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	requestBill := mb.CreateBillRequest{
		Amount:         mb.Money(order.Amount / 100.0),
		CustomOrderId:  strconv.FormatUint(order.Id, 10),
		Description:    fmt.Sprintf("Заказ #%d", order.Id),
		SendLetter:     false,
		Testing:        false,
		ReceiptContact: customer.Email.String,
		ReceiptItems:   nil,
		UnixTimestamp:  time.Now().Unix(),
		Salt:           gofakeit.UUID(),
	}

	receiptItems := make([]mb.ReceiptItem, 0)
	for _, i := range order.ReceiptItems {
		ri := mb.ReceiptItem{
			Name:          i.Name,
			Quantity:      i.Quantity,
			Price:         mb.Money(i.Price / 100.0),
			SNO:           mb.SNOUsnIncome,
			PaymentObject: mb.PaymentObjectService,
			PaymentMethod: mb.PaymentMethodFullPrepayment,
			VAT:           mb.VATNone,
		}
		receiptItems = append(receiptItems, ri)
	}
	requestBill.ReceiptItems = receiptItems

	opts := mb.MerchantOptions{
		Merchant:  user.MerchantId,
		SecretKey: user.SecretKey,
	}
	bill, err := mb.CreateBill(ctx, requestBill, opts, http.DefaultClient)
	if err != nil {
		return fmt.Errorf("failed to create bill: %w", err)
	}

	err = b.storage.UpdatePaymentLink(ctx, paymentId, bill.Url, bill.Id)
	if err != nil {
		return fmt.Errorf("failed to generate payment link: %w", err)
	}
	return nil
}
