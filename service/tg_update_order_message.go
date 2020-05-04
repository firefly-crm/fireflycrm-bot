package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (s Service) updateOrderMessage(ctx context.Context, bot *tg.BotAPI, orderId uint64, flowCompleted bool) error {
	order, err := s.OrderBook.GetOrder(ctx, orderId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	var customer *types.Customer
	if order.CustomerId.Valid {
		c, err := s.Users.GetCustomer(ctx, uint64(order.CustomerId.Int64))
		if err != nil {
			logrus.Error("failed to get customer: %w", err)
		}
		customer = &c
	}

	chatId := int64(order.UserId)
	messageId := int(order.MessageId)

	editMessage := tg.NewEditMessageText(chatId, messageId, order.MessageString(customer))
	editMessage.ParseMode = "markdown"
	editMessage.DisableWebPagePreview = true
	var markup tg.InlineKeyboardMarkup
	if flowCompleted {
		markup = startOrderInlineKeyboard()
	} else {
		markup = cancelInlineKeyboard()
	}
	editMessage.ReplyMarkup = &markup

	_, err = bot.Send(editMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
