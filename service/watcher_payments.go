package service

import (
	"context"
	"fmt"
	mb "github.com/DarthRamone/modulbank-go"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s Service) startPaymentsWatcher(ctx context.Context, bot *tg.BotAPI) {
	ticker := time.NewTicker(time.Minute * 1)
	defer func() {
		ticker.Stop()
	}()

outsideLoop:
	for {
		select {
		case <-ctx.Done():
			break outsideLoop
		case <-ticker.C:
			err := s.checkPayments(ctx, bot)
			if err != nil {
				logrus.Errorf("failed to check payments: %v", err.Error())
			}
		}
	}
}

func (s Service) checkPayments(ctx context.Context, bot *tg.BotAPI) error {
	payments, err := s.OrderBook.GetBankPayments(ctx)
	if err != nil {
		return fmt.Errorf("failed to get payments: %w", err)
	}

	for _, p := range payments {
		order, err := s.OrderBook.GetOrder(ctx, p.OrderId)
		if err != nil {
			logrus.Errorf("failed to get order: %w", err)
			continue
		}

		user, err := s.Storage.GetUser(ctx, order.UserId)
		if err != nil {
			logrus.Errorf("failed to get user: %w", err)
			continue
		}

		opts := mb.MerchantOptions{
			Merchant:  user.MerchantId,
			SecretKey: user.SecretKey,
		}

		logrus.Debugf("bank payment id: %s\n", p.BankPaymentId)

		bill, err := mb.GetBill(ctx, p.BankPaymentId, opts, http.DefaultClient)
		if err != nil {
			logrus.Errorf("failed to get bill: %w", err)
			continue
		}

		if bill.Paid == 1 {
			err := s.Storage.SetPaymentPaid(ctx, p.Id)
			if err != nil {
				logrus.Errorf("failed to set payment paid: %w", err)
				continue
			}

			messages, err := s.Storage.GetMessagesForOrder(ctx, order.Id)
			if err != nil {
				logrus.Errorf("failed to get messages for order: %w", err)
				continue
			}

			if len(messages) == 0 {
				logrus.Errorf("no messages for order found")
				continue
			}

			msg := tg.NewMessage(int64(user.Id), "Заказ оплачен")
			msg.ReplyToMessageID = int(messages[0].Id)
			msg.ReplyMarkup = notifyReadInlineKeyboard()
			_, err = bot.Send(msg)
			if err != nil {
				logrus.Errorf("failed to send message to chat: %v", err)
				continue
			}

			err = s.updateOrderMessage(ctx, bot, messages[0].Id, true)
			if err != nil {
				logrus.Errorf("failed to update order message: %v", err)
				continue
			}
		} else if bill.Expired == 1 {
			err := s.Storage.SetPaymentExpired(ctx, p.Id)
			if err != nil {
				logrus.Errorf("failed to set payment expired: %w", err)
				continue
			}
		}
	}

	return nil
}
