package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processCustomerEditEmail(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery) error {
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

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

	err = s.OrderBook.UpdateOrderState(ctx, order.Id, types.WaitingCustomerEmail)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
