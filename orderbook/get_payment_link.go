package orderbook

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
)

func (b orderBook) GeneratePaymentLink(ctx context.Context, paymentId uint64) error {
	url := gofakeit.URL()
	err := b.storage.UpdatePaymentLink(ctx, paymentId, url)
	if err != nil {
		return fmt.Errorf("failed to generate payment link: %w", err)
	}
	return nil
}
