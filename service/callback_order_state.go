package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/DarthRamone/telegram-bot-api"
)

func (s Service) processOrderStateCallback(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, state types.OrderState) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderState(ctx, order.Id, state)
	if err != nil {
		return fmt.Errorf("failed to update order state(%s): %w", state, err)
	}

	if state == types.OrderStateDeleted {
		err = s.OrderBook.UpdateOrderMessageDisplayMode(ctx, messageId, types.DisplayModeDeleted)
		if err != nil {
			return fmt.Errorf("failed to update display mode: %w", err)
		}
	}

	if state == types.OrderStateForming {
		orderMessage, err := s.OrderBook.GetOrderMessage(ctx, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order message: %w", err)
		}

		if orderMessage.DisplayMode == types.DisplayModeDeleted {
			err = s.OrderBook.UpdateOrderMessageDisplayMode(ctx, messageId, types.DisplayModeFull)
			if err != nil {
				return fmt.Errorf("failed to update display mode: %w", err)
			}
		}
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}

func (s Service) processOrderEditStateCallback(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, state types.EditState) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, state)
	if err != nil {
		return fmt.Errorf("failed to update order state(%s): %w", state, err)
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
