package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (s Service) processAddPaymentCallback(ctx context.Context, cbq *tg.CallbackQuery, method types.PaymentMethod) error {
	messageId := cbq.Message.MessageID

	order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
	if err != nil {
		return fmt.Errorf("failed to get order by message id: %w", err)
	}

	_, err = s.OrderBook.AddPayment(ctx, order.Id, method)
	if err != nil {
		return fmt.Errorf("failed to add payment to order: %w", err)
	}

	return nil
}

func (s Service) processPartialPaymentCallback(ctx context.Context, bot *tg.BotAPI, cbq *tg.CallbackQuery) error {
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

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingPaymentAmount)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}

//if amount is 0 then full payment
func (s Service) processPaymentCallback(ctx context.Context, bot *tg.BotAPI, messageId uint64, amount uint32) error {
	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	if !order.ActivePaymentId.Valid {
		return fmt.Errorf("active bill id is nil")
	}

	paymentId := uint64(order.ActivePaymentId.Int64)
	if amount == 0 {
		amount = order.Amount - order.PayedAmount
	}

	var payment types.Payment
	for _, p := range order.Payments {
		if p.Id == paymentId {
			payment = p
			break
		}
	}

	defer func() {
		if err := s.deleteHint(ctx, bot, order); err != nil {
			logrus.Error("failed to delete hint: %v", err.Error())
		}
	}()

	err = s.OrderBook.UpdatePaymentAmount(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to update payment amount: %w", err)
	}

	if payment.PaymentMethod == types.PaymentMethodAcquiring {
		err := s.OrderBook.GeneratePaymentLink(ctx, paymentId)
		if err != nil {
			return fmt.Errorf("failed to generate payment link: %w", err)
		}
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
