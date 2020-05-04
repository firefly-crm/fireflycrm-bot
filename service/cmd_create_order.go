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

	messageText := fmt.Sprintf(`
*Заказ №%d*

*Позиции:*

*Итого: 0.00р*

*Данные клиента:*
`, orderId)

	messageReplyMarkup := startOrderInlineKeyboard()

	deleteMessage := tg.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err = bot.DeleteMessage(deleteMessage)
	if err != nil {
		log.Printf("failed to delete command message: %w", err)
	}

	msg := tg.NewMessage(update.Message.Chat.ID, messageText)
	msg.ReplyMarkup = messageReplyMarkup
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

	return nil
}
