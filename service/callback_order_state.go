package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processOrderStateCallback(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, state types.OrderState) error {
	messageId := callbackQuery.Message.MessageID

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderState(ctx, order.Id, state)
	if err != nil {
		return fmt.Errorf("failed to update order state(%s): %w", state, err)
	}

	err = s.updateOrderMessage(ctx, bot, order.Id, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
