package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processPaymentRemove(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, paymentId uint64) error {
	messageId := callbackQuery.Message.MessageID

	err := s.OrderBook.RemovePayment(ctx, paymentId)
	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, order.Id, true)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}
