package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/DarthRamone/telegram-bot-api"
)

func (s Service) processCustomerEditInstagram(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery) error {
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	err := s.Users.SetActiveOrderMessageForUser(ctx, uint64(chatId), uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to set active order for user: %w", err)
	}

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(chatId, replyEnterCustomerInstagram)
	hint, err := bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingCustomerInstagram)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processCustomerEditEmail(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery) error {
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	err := s.Users.SetActiveOrderMessageForUser(ctx, uint64(chatId), uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to set active order for user: %w", err)
	}

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(chatId, replyEnterCustomerEmail)
	hint, err := bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingCustomerEmail)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processCustomerEditPhone(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery) error {
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	err := s.Users.SetActiveOrderMessageForUser(ctx, uint64(chatId), uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to set active order for user: %w", err)
	}

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(chatId, replyEnterCustomerPhone)
	hint, err := bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingCustomerPhone)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
