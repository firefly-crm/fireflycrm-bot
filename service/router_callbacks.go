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

func (s Service) processCallback(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	callbackQuery := update.CallbackQuery
	chatId := callbackQuery.Message.Chat.ID
	messageId := callbackQuery.Message.MessageID

	var paymentMethod types.PaymentMethod

	var markup tg.InlineKeyboardMarkup
	callbackData := callbackQuery.Data

	logrus.Println(callbackData)

	deleteOrderMessage := false

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
	case kbDataRemoveItem:
		var err error
		markup, err = itemsListInlineKeyboard(ctx, s, uint64(messageId), "remove")
		if err != nil {
			return fmt.Errorf("failed to get markup for remove items list: %w", err)
		}
		break
	case kbDataEditItem:
		var err error
		markup, err = itemsListInlineKeyboard(ctx, s, uint64(messageId), "edit")
		if err != nil {
			return fmt.Errorf("failed to get markup for edit items list: %w", err)
		}
		break
	case kbDataCustomer:
		var err error
		markup, err = customerInlineKeyboard(ctx, s, uint64(messageId))
		if err != nil {
			return fmt.Errorf("failed to get markup for customer action: %w", err)
		}
		break
	case kbDataPayment:
		markup = paymentInlineKeyboard()
		break
	case kbDataPaymentCard:
		paymentMethod = types.Card2Card
		markup = paymentAmountInlineKeyboard()
		err := s.processAddPaymentCallback(ctx, callbackQuery, paymentMethod)
		if err != nil {
			return fmt.Errorf("failed to add card payment to order: %w", err)
		}
		break
	case kbDataPaymentCash:
		paymentMethod = types.Cash
		markup = paymentAmountInlineKeyboard()
		err := s.processAddPaymentCallback(ctx, callbackQuery, paymentMethod)
		if err != nil {
			return fmt.Errorf("failed to add cash payment to order: %w", err)
		}
		break
	case kbDataPaymentLink:
		paymentMethod = types.Acquiring
		markup = paymentAmountInlineKeyboard()
		err := s.processAddPaymentCallback(ctx, callbackQuery, paymentMethod)
		if err != nil {
			return fmt.Errorf("failed to add link payment to order: %w", err)
		}
		break
	case kbDataFullPayment:
		markup = startOrderInlineKeyboard()
		err := s.processPaymentCallback(ctx, bot, uint64(messageId), 0)
		if err != nil {
			return fmt.Errorf("failed to process full payment callback: %w", err)
		}
		break
	case kbDataPartialPayment:
		markup = cancelInlineKeyboard()
		err := s.processPartialPaymentCallback(ctx, bot, callbackQuery)
		if err != nil {
			return fmt.Errorf("failed to process partial payment callback: %w", err)
		}
		break
	case kbDataRefundPayment:
		var err error
		markup, err = paymentsListInlineKeyboard(ctx, s, uint64(messageId), "refund")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case kbDataPartialRefund:
		markup = cancelInlineKeyboard()
		err := s.processPartialRefundCallback(ctx, bot, callbackQuery)
		if err != nil {
			return fmt.Errorf("failed to process partial refund callback: %w", err)
		}
	case kbDataFullRefund:
		markup = startOrderInlineKeyboard()
		err := s.processRefundCallback(ctx, bot, uint64(messageId), 0)
		if err != nil {
			return fmt.Errorf("failed to process refund callback: %w", err)
		}
	case kbDataRemovePayment:
		var err error
		markup, err = paymentsListInlineKeyboard(ctx, s, uint64(messageId), "remove")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case kbDataOrderActions:
		var err error
		markup, err = orderActionsInlineKeyboard(ctx, s, uint64(messageId))
		if err != nil {
			return fmt.Errorf("failed to get order actions markup: %w", err)
		}
	case kbDataOrderDone:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.Completed)
		if err != nil {
			return fmt.Errorf("failed to process order done callback: %w", err)
		}
		markup = startOrderInlineKeyboard()
	case kbDataOrderRestore:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.StandBy)
		if err != nil {
			return fmt.Errorf("failed to process order done callback: %w", err)
		}
		markup = startOrderInlineKeyboard()
	case kbDataOrderDelete:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.Completed)
		if err != nil {
			return fmt.Errorf("failed to process order done callback: %w", err)
		}
		deleteOrderMessage = true
	default:
		args := strings.Split(callbackData, "_")
		entity := args[0]
		action := args[1]

		argsCount := len(args)

		if entity == "customer" {
			switch args[2] {
			case "name":
			case "email":
				markup = cancelInlineKeyboard()
				err := s.processCustomerEditEmail(ctx, bot, callbackQuery)
				if err != nil {
					return fmt.Errorf("failed to process item edit name callback: %w", err)
				}
			case "phone":
			}
		}

		if entity == "payment" {
			strId := args[len(args)-1]
			id, err := strconv.ParseUint(strId, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse id: %w", err)
			}

			switch args[1] {
			case "remove":
				err := s.processPaymentRemove(ctx, bot, callbackQuery, id)
				if err != nil {
					return fmt.Errorf("failed to process remove payment callback: %w", err)
				}
				markup = startOrderInlineKeyboard()
				break
			case "refund":
				order, err := s.OrderBook.GetOrderByMessageId(ctx, uint64(messageId))
				if err != nil {
					return fmt.Errorf("failed to get order: %w", err)
				}
				err = s.OrderBook.SetActivePaymentId(ctx, order.Id, id)
				if err != nil {
					return fmt.Errorf("failed to set active refund payment: %w", err)
				}
				markup = refundAmountInlineKeyboard()
				break
			}
		}

		if entity == "item" {
			if argsCount == 3 {
				id, err := strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse id: %w", err)
				}

				if action == "edit" {
					markup = editItemActionsInlineKeyboard(id)
				}

				if action == "remove" {
					err := s.processItemRemove(ctx, bot, callbackQuery, id)
					if err != nil {
						return fmt.Errorf("failed to remove item: %w", err)
					}
				}
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

	var msg tg.Chattable

	if deleteOrderMessage {
		msg = tg.NewDeleteMessage(chatId, messageId)
	} else {
		msg = tg.NewEditMessageReplyMarkup(chatId, messageId, markup)
	}

	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to update order message: %w", err)
	}

	return nil
}
