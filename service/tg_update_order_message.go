package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (s Service) updateOrderMessage(ctx context.Context, bot *tg.BotAPI, messageId uint64, flowCompleted bool) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	logrus.Info("update order message; got order")

	orderMessage, err := s.OrderBook.GetOrderMessage(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order message: %w", err)
	}
	logrus.Info("update order message; got order message")

	var customer *types.Customer
	if order.CustomerId.Valid {
		logrus.Info("update order message; start get order")
		c, err := s.Users.GetCustomer(ctx, uint64(order.CustomerId.Int64))
		if err != nil {
			logrus.Error("failed to get customer: %w", err)
		}
		customer = &c
		logrus.Info("update order message; got customer")
	}

	chatId := int64(order.UserId)

	editMessage := tg.NewEditMessageText(chatId, int(messageId), order.MessageString(customer, orderMessage.DisplayMode))
	//editMessage.ParseMode = "markdown"
	editMessage.DisableWebPagePreview = true
	var markup tg.InlineKeyboardMarkup
	if flowCompleted {
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	} else {
		markup = cancelInlineKeyboard()
	}
	editMessage.ReplyMarkup = &markup

	logrus.Infof("update order message(%d); send message:\n\n%s\n\n", messageId, editMessage.Text)
	_, err = bot.Send(editMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	logrus.Info("update order message; message sent")

	return nil
}

/*
*Заказ #22* _(Формируется)_
*Создан:* 15.05.2020
*Сумма:* 5000.00₽
*Оплачен:* полностью

*Позиции*
`- Доставка 5000.00₽ x1`

*Клиент*
*E-Mail:* molekyla89@mail.ru

*Данные по оплате*

*Платеж #1.* Перевод на карту.
*Сумма:* 5000.00₽
*Оплачен:* 15 May 2020 03:31
*/
