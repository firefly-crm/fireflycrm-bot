package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processItemRemove(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, itemId uint64) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	err := s.OrderBook.RemoveItem(ctx, itemId)
	if err != nil {
		return fmt.Errorf("failed to remove receipt item: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}
