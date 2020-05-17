package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processOrderDisplayModeCallback(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, mode types.DisplayMode) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	err := s.OrderBook.UpdateOrderMessageDisplayMode(ctx, messageId, mode)
	if err != nil {
		return fmt.Errorf("failed to update display mode: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
