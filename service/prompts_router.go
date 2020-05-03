package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (s Service) deleteHint(ctx context.Context, bot *tg.BotAPI, order types.Order) error {
	if !order.HintMessageId.Valid {
		return fmt.Errorf("hint id is nil")
	}

	deleteMessage := tg.NewDeleteMessage(int64(order.UserId), int(order.HintMessageId.Int64))
	_, err := bot.Send(deleteMessage)
	if err != nil {
		return fmt.Errorf("failed to delete hind: %w", err)
	}

	markup := startOrderInlineKeyboard()
	editMarkup := tg.NewEditMessageReplyMarkup(int64(order.UserId), int(order.HintMessageId.Int64), markup)
	_, err = bot.Send(editMarkup)
	if err != nil {
		return fmt.Errorf("failed to send new markup: %w", err)
	}

	return nil
}

func (s Service) updateOrderMessage(ctx context.Context, bot *tg.BotAPI, orderId uint64, flowCompleted bool) error {
	order, err := s.OrderBook.GetOrder(ctx, orderId)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}

	chatId := int64(order.UserId)
	messageId := int(order.MessageId)

	editMessage := tg.NewEditMessageText(chatId, messageId, order.MessageString())
	editMessage.ParseMode = "markdown"
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

func (s Service) processPrompt(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	userId := uint64(update.Message.From.ID)
	activeOrder, err := s.OrderBook.GetActiveOrderForUser(ctx, userId)
	if err != nil {
		return fmt.Errorf("failed to get active order for user: %w", err)
	}

	flowCompleted := true

	defer func() {
		delMessage := tg.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
		_, err := bot.Send(delMessage)
		if err != nil {
			logrus.Errorf("failed to delete message: %v", err)
		}

		err = s.updateOrderMessage(ctx, bot, activeOrder.Id, flowCompleted)
		if err != nil {
			logrus.Errorf("failed to update order message: %w", err)
		}
	}()

	text := strings.TrimSpace(update.Message.Text)

	switch activeOrder.State {
	case types.WaitingItemName:
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
		} else {
			err := s.deleteHint(ctx, bot, activeOrder)
			if err != nil {
				return fmt.Errorf("failed to remove hint: %w", err)
			}
		}

		break
	case types.WaitingItemPrice:
		if !activeOrder.ActiveItemId.Valid {
			return fmt.Errorf("active order item id doesnt exists")
		}

		price, err := strconv.Atoi(text)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}

		receiptItemId := uint64(activeOrder.ActiveItemId.Int64)
		err = s.OrderBook.UpdateReceiptItemPrice(ctx, uint32(price*100), receiptItemId)
		if err != nil {
			return fmt.Errorf("failed to change item price: %w", err)
		}

		err = s.deleteHint(ctx, bot, activeOrder)
		if err != nil {
			return fmt.Errorf("failed to remove hint: %w", err)
		}

		break
	case types.WaitingItemQuantity:
		break
	}

	return nil
}
