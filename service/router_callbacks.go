package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (s Service) processCallback(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	callbackQuery := update.CallbackQuery
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	var markup tg.InlineKeyboardMarkup
	callbackData := callbackQuery.Data

	logrus.Println(callbackData)

	switch callbackData {
	case kbDataItems:
		markup = orderItemsInlineKeyboard()
		break
	case kbDataBack:
		markup = startOrderInlineKeyboard()
		break
	case kbDataCancel:
		markup = startOrderInlineKeyboard()
		err := s.processCancelCallback(ctx, bot, callbackQuery)
		if err != nil {
			return fmt.Errorf("failed to process cancel callback: %w", err)
		}
		break
	case kbDataAddItem:
		markup = cancelInlineKeyboard()
		err := s.processAddItemCallack(ctx, bot, callbackQuery)
		if err != nil {
			return fmt.Errorf("failed to process add item callback: %w", err)
		}
		break
	case kbDataEditItem:
		var err error
		markup, err = itemsListInlineKeyboard(ctx, s, uint64(messageId))
		if err != nil {
			return fmt.Errorf("failed to get markup for items list: %w", err)
		}
		break
	default:
		args := strings.Split(callbackData, "_")
		entity := args[0]
		//action := args[1]

		argsCount := len(args)

		if entity == "item" {
			if argsCount == 3 {
				id, err := strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse id: %w", err)
				}
				markup = editItemActionsInlineKeyboard(id)
			}

			if argsCount == 4 {
				id, err := strconv.ParseUint(args[3], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse id: %w", err)
				}

				switch args[2] {
				case "qty":
					markup = cancelInlineKeyboard()
					err = s.processItemEditQty(ctx, bot, callbackQuery, id)
					if err != nil {
						return fmt.Errorf("failed to process item edit qty callback: %w", err)
					}
				case "name":
					markup = cancelInlineKeyboard()
					err = s.processItemEditName(ctx, bot, callbackQuery, id)
					if err != nil {
						return fmt.Errorf("failed to process item edit name callback: %w", err)
					}
				case "price":
					markup = cancelInlineKeyboard()
					err = s.processItemEditPrice(ctx, bot, callbackQuery, id)
					if err != nil {
						return fmt.Errorf("failed to process item edit price callback: %w", err)
					}
				}
			}
		}
	}

	edit := tg.NewEditMessageReplyMarkup(chatId, messageId, markup)

	_, err := bot.Send(edit)
	if err != nil {
		return fmt.Errorf("failed to update order message inline keyboard: %w", err)
	}

	return nil
}
