package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (s Service) createOrder(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	userId := uint64(update.Message.From.ID)

	orderId, err := s.OrderBook.CreateOrder(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	messageText := fmt.Sprintf(`*Заказ №%d*`, orderId)

	deleteMessage := tg.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err = bot.DeleteMessage(deleteMessage)
	if err != nil {
		log.Printf("failed to delete command message: %w", err)
	}

	msg := tg.NewMessage(update.Message.Chat.ID, messageText)
	msg.ParseMode = "markdown"

	var orderMessage tg.Message
	orderMessage, err = bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateMessageForOrder(ctx, orderId, uint64(orderMessage.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update message id for order: %v", err)
	}

	messageReplyMarkup, err := startOrderInlineKeyboard(ctx, s, uint64(orderMessage.MessageID))
	if err != nil {
		return fmt.Errorf("failed to get start order inline markup: %w", err)
	}

	editMessage := tg.NewEditMessageReplyMarkup(update.Message.Chat.ID, orderMessage.MessageID, messageReplyMarkup)
	_, err = bot.Send(editMessage)
	if err != nil {
		return fmt.Errorf("failed to set new markup: %w", err)
	}

	return nil
}
