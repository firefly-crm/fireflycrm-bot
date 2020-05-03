package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

func (s Service) processCallback(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	callbackQuery := update.CallbackQuery
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID
	var markup tg.InlineKeyboardMarkup

	log.Printf("callback data: %s\n", callbackQuery.Data)

	switch callbackQuery.Data {
	case kbDataItems:
		markup = orderItemsInlineKeyboard()
		break
	case kbDataBack:
		markup = startOrderInlineKeyboard()
		break
	case kbDataCancel:
		markup = startOrderInlineKeyboard()
		err := s.processCancelCallback(ctx, bot, uint64(messageId))
		if err != nil {
			return fmt.Errorf("failed to process cancel callback: %w", err)
		}

		break
	case kbDataAddItem:
		markup = cancelInlineKeyboard()

		order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
		if err != nil {
			return fmt.Errorf("failed to get order by message id: %w", err)
		}
		hintMessage := tg.NewMessage(chatId, replyEnterItemName)
		hint, err := bot.Send(hintMessage)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		err = s.Users.SetActiveOrderForUser(ctx, order.UserId, order.Id)
		if err != nil {
			return fmt.Errorf("failed to set active order for user: %w", err)
		}

		_, err = s.OrderBook.AddItem(ctx, order.Id)
		if err != nil {
			return fmt.Errorf("failed to add item to order")
		}

		err = s.OrderBook.UpdateHintMessageForOrder(ctx, order.Id, uint64(hint.MessageID))
		if err != nil {
			return fmt.Errorf("failed to update hint message: %w", err)
		}

		err = s.OrderBook.UpdateOrderState(ctx, order.Id, types.WaitingItemName)
		if err != nil {
			return fmt.Errorf("failed to update order state: %w", err)
		}

		break
	}

	edit := tg.NewEditMessageReplyMarkup(chatId, messageId, markup)

	_, err := bot.Send(edit)
	if err != nil {
		return fmt.Errorf("failed to update order message inline keyboard: %w", err)
	}

	return nil
}
