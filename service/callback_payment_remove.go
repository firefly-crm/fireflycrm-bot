package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processPaymentRemove(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, paymentId uint64) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	err := s.OrderBook.RemovePayment(ctx, paymentId)
	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}
