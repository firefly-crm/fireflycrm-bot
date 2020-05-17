package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/DarthRamone/telegram-bot-api"
)

func (s Service) processAddKnownItem(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, data string) error {
	chatId := callbackQuery.Message.Chat.ID
	messageId := uint64(callbackQuery.Message.MessageID)

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	hintMessage := tg.NewMessage(chatId, replyEnterItemPrice)
	hint, err := bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	var itemId uint64
	itemId, err = s.OrderBook.AddItem(ctx, order.Id)
	if err != nil {
		return fmt.Errorf("failed to add item to order")
	}

	name := "Unknown item"
	switch data {
	case kbDataDelivery:
		name = "Доставка"
	case kbDataLingerieSet:
		name = "Комплект нижнего белья"
	}

	err = s.OrderBook.UpdateReceiptItemName(ctx, name, uint64(chatId), itemId)
	if err != nil {
		return fmt.Errorf("failed to set delivery name: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemPrice)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

func (s Service) processAddItemCallack(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery) error {
	chatId := callbackQuery.Message.Chat.ID
	messageId := uint64(callbackQuery.Message.MessageID)

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(chatId, replyEnterItemName)
	hint, err := bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	_, err = s.OrderBook.AddItem(ctx, order.Id)
	if err != nil {
		return fmt.Errorf("failed to add item to order")
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemName)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
