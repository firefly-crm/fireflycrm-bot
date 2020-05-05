package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/badoux/checkmail"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (s Service) processPrompt(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	userId := uint64(update.Message.From.ID)
	activeMessageId, err := s.OrderBook.GetActiveOrderMessageIdForUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get active message id: %w", err)
	}
	activeOrder, err := s.OrderBook.GetActiveOrderForUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get active order for user: %w", err)
	}

	deleteHint := true
	standBy := true
	flowCompleted := true

	defer func() {
		if deleteHint {
			err = s.deleteHint(ctx, bot, activeOrder)
			if err != nil {
				logrus.Errorf("failed to remove hint: %v", err)
			}
		}

		if standBy {
			err = s.OrderBook.UpdateOrderEditState(ctx, activeOrder.Id, types.EditStateNone)
			if err != nil {
				logrus.Errorf("failed to set standby mode: %w", err)
			}
		}

		delMessage := tg.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		_, err := bot.Send(delMessage)
		if err != nil {
			logrus.Errorf("failed to delete message: %v", err)
		}

		err = s.updateOrderMessage(ctx, bot, activeMessageId, flowCompleted)
		if err != nil {
			logrus.Errorf("failed to update order message: %w", err)
		}
	}()

	text := strings.TrimSpace(update.Message.Text)

	switch activeOrder.EditState {
	case types.EditStateWaitingItemName:
		if !activeOrder.ActiveItemId.Valid {
			return fmt.Errorf("active order item id doesnt exists")
		}

		receiptItemId := uint64(activeOrder.ActiveItemId.Int64)

		err := s.OrderBook.UpdateReceiptItemName(ctx, text, userId, receiptItemId)
		if err != nil {
			return fmt.Errorf("failed to change item name: %w", err)
		}

		item, err := s.OrderBook.GetReceiptItem(ctx, receiptItemId)
		if err != nil {
			return fmt.Errorf("failed to get receipt item: %w", err)
		}

		if !item.Initialised {
			err := s.setWaitingForPrice(ctx, bot, activeOrder)
			if err != nil {
				return fmt.Errorf("failed to change order state: %w", err)
			}
			flowCompleted = false
			deleteHint = false
			standBy = false
		}

		break
	case types.EditStateWaitingItemPrice:
		if !activeOrder.ActiveItemId.Valid {
			return fmt.Errorf("active order item id doesnt exists")
		}

		text = strings.Trim(text, "₽р$РP")

		price, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}

		receiptItemId := uint64(activeOrder.ActiveItemId.Int64)
		err = s.OrderBook.UpdateReceiptItemPrice(ctx, uint32(price*100), receiptItemId)
		if err != nil {
			return fmt.Errorf("failed to change item price: %w", err)
		}

		break
	case types.EditStateWaitingItemQuantity:
		if !activeOrder.ActiveItemId.Valid {
			return fmt.Errorf("active order item id doesnt exists")
		}

		qty, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}

		receiptItemId := uint64(activeOrder.ActiveItemId.Int64)
		err = s.OrderBook.UpdateReceiptItemQty(ctx, qty, receiptItemId)
		if err != nil {
			return fmt.Errorf("failed to change item quantity: %w", err)
		}

		break
	case types.EditStateWaitingCustomerEmail:
		err = checkmail.ValidateFormat(text)
		if err != nil {
			return fmt.Errorf("email validation failed: %w", err)
		}

		_, err = s.OrderBook.UpdateCustomerEmail(ctx, text, activeOrder.Id)
		if err != nil {
			return fmt.Errorf("failed to update customer email: %w", err)
		}

		break
	case types.EditStateWaitingPaymentAmount:
		if !activeOrder.ActivePaymentId.Valid {
			return fmt.Errorf("active payment id doesnt exists")
		}

		amount, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %w", err)
		}

		err = s.processPaymentCallback(ctx, bot, activeMessageId, uint32(amount*100))
		if err != nil {
			return fmt.Errorf("failed to proces payment callback")
		}

		break
	case types.EditStateWaitingRefundAmount:
		if !activeOrder.ActivePaymentId.Valid {
			return fmt.Errorf("active payment id doesnt exists")
		}

		amount, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %w", err)
		}

		err = s.processRefundCallback(ctx, bot, activeOrder, activeMessageId, uint32(amount*100))
		if err != nil {
			return fmt.Errorf("failed to proces refund callback")
		}

		break
	case types.EditStateWaitingCustomerInstagram:
		text = strings.Trim(text, "@")

		_, err = s.OrderBook.UpdateCustomerInstagram(ctx, text, activeOrder.Id)
		if err != nil {
			return fmt.Errorf("failed to update customer email: %w", err)
		}

		break
	}

	return nil
}
