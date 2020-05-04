package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (s Service) processPaymentRefund(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, paymentId uint64, amount uint32) error {
	messageId := callbackQuery.Message.MessageID

	err := s.OrderBook.RefundPayment(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, order.Id, true)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}

func (s Service) processRefundCallback(ctx context.Context, bot *tg.BotAPI, messageId uint64, amount uint32) error {

	//TODO: Refund payment at ModulBank

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	if !order.ActivePaymentId.Valid {
		return fmt.Errorf("active payment id is nil")
	}

	paymentId := uint64(order.ActivePaymentId.Int64)
	if amount == 0 {
		for _, p := range order.Payments {
			if p.Id == paymentId {
				amount = p.Amount
			}
		}
	}

	defer func() {
		if err := s.deleteHint(ctx, bot, order); err != nil {
			logrus.Error("failed to delete hint: %v", err)
		}
	}()

	err = s.OrderBook.RefundPayment(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, order.Id, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}

func (s Service) processPartialRefundCallback(ctx context.Context, bot *tg.BotAPI, cbq *tg.CallbackQuery) error {
	chatId := cbq.Message.Chat.ID
	messageId := cbq.Message.MessageID

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}
	hintMessage := tg.NewMessage(chatId, replyEnterAmount)
	hint, err := bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
	if err != nil {
		return fmt.Errorf("failed to update hint message: %w", err)
	}

	err = s.OrderBook.UpdateOrderState(ctx, order.Id, types.WaitingRefundAmount)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
