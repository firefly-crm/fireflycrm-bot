package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/badoux/checkmail"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

func (s Service) processPrompt(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	userId := uint64(update.Message.From.ID)
	activeMessageId, err := s.OrderBook.GetActiveOrderMessageIdForUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get active message id: %w", err)
	}
	logrus.Info("got active message id")

	activeOrder, err := s.OrderBook.GetActiveOrderForUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get active order for user: %w", err)
	}
	logrus.Info("got active order")

	deleteHint := true
	standBy := true
	flowCompleted := true

	defer func() {
		if deleteHint {
			logrus.Info("prompts; delete hint")
			err = s.deleteHint(ctx, bot, activeOrder)
			if err != nil {
				logrus.Errorf("failed to remove hint: %v", err)
			}
		}

		if standBy {
			logrus.Info("prompts; edit state; stand by")
			err = s.OrderBook.UpdateOrderEditState(ctx, activeOrder.Id, types.EditStateNone)
			if err != nil {
				logrus.Errorf("failed to set standby mode: %w", err)
			}
		}

		logrus.Info("prompts; deleting message")
		delMessage := tg.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		_, err := bot.Send(delMessage)
		if err != nil {
			logrus.Errorf("failed to delete message: %v", err)
		}

		logrus.Info("prompts; updating message")
		err = s.updateOrderMessage(ctx, bot, activeMessageId, flowCompleted)
		if err != nil {
			logrus.Errorf("failed to update order message: %v", err)
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
		logrus.Info("updated receipt item name")

		item, err := s.OrderBook.GetReceiptItem(ctx, receiptItemId)
		if err != nil {
			return fmt.Errorf("failed to get receipt item: %w", err)
		}
		logrus.Info("got receipt item")

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

		price, err := parseAmount(text)
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

		qty, err := parseAmount(text)
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

		amount, err := parseAmount(text)
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

		amount, err := parseAmount(text)
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
	case types.EditStateWaitingCustomerPhone:
		text = strings.Trim(text, "+")

		phone, err := parsePhone(text)
		if err != nil {
			return fmt.Errorf("failed to parse phone number: %w", err)
		}

		_, err = s.Storage.UpdateCustomerPhone(ctx, phone, activeOrder.Id)
		if err != nil {
			return fmt.Errorf("failed to update customer phone: %w", err)
		}

		break
	}

	return nil
}

func parseAmount(text string) (int, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return 0, fmt.Errorf("failed to parse text: %w", err)
	}
	numericStr := reg.ReplaceAllString(text, "")
	res, err := strconv.ParseInt(numericStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse")
	}
	return int(res), nil
}

func parsePhone(text string) (string, error) {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return "", fmt.Errorf("failed to parse text: %w", err)
	}
	numericStr := reg.ReplaceAllString(text, "")

	if len(numericStr) == 10 {
		numericStr = "7" + numericStr
	}

	if numericStr[0] == '8' {
		numericStr = "7" + numericStr[1:]
	}

	digits := strings.Split(numericStr, "")

	phone := fmt.Sprintf("+%s(%s)%s-%s-%s",
		digits[0],
		strings.Join(digits[1:4], ""),
		strings.Join(digits[4:7], ""),
		strings.Join(digits[7:9], ""),
		strings.Join(digits[9:], ""))

	return phone, nil
}
