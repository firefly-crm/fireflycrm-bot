package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/common/logger"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) processPaymentRefund(ctx context.Context, bot *tg.BotAPI, callbackQuery *tg.CallbackQuery, paymentId uint64, amount uint32) error {
	messageId := uint64(callbackQuery.Message.MessageID)

	err := s.OrderBook.RefundPayment(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
	if err != nil {
		return fmt.Errorf("failed to refresh order message: %w", err)
	}

	return nil
}

func (s Service) processRefundCallback(ctx context.Context, bot *tg.BotAPI, order types.Order, messageId uint64, amount uint32) error {
	log := logger.FromContext(ctx)

	//TODO: Refund payment at ModulBank

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
			log.Errorf("failed to delete hint: %v", err.Error())
		}
	}()

	err := s.OrderBook.RefundPayment(ctx, paymentId, amount)
	if err != nil {
		return fmt.Errorf("failed to refund payment: %w", err)
	}

	err = s.updateOrderMessage(ctx, bot, messageId, true)
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

	err = s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingRefundAmount)
	if err != nil {
		return fmt.Errorf("failed to update order state: %w", err)
	}

	return nil
}
