package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processCancelCallback(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	err = s.OrderBook.UpdateOrderState(ctx, order.Id, types.StandBy)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	if order.ActiveItemId.Valid {
		err = s.OrderBook.RemoveItem(ctx, uint64(order.ActiveItemId.Int64))
		if err != nil {
			return fmt.Errorf("failed to delete receipt item: %w", err)
		}
	}

	err = s.updateOrderMessage(ctx, bot, order.Id, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	err = s.deleteHint(ctx, bot, order)
	if err != nil {
		return fmt.Errorf("failed to delete hint: %w", err)
	}

	return nil
}
