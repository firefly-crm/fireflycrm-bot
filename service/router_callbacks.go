package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func (s Service) processCallback(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	callbackQuery := update.CallbackQuery
	chatId := callbackQuery.Message.Chat.ID
	messageId := uint64(callbackQuery.Message.MessageID)

	var markup tg.InlineKeyboardMarkup
	callbackData := callbackQuery.Data

	shouldDelete := false
	logrus.Infof("processing callback data: %s", callbackData)

	switch callbackData {
	case kbDataItems:
		markup = orderItemsInlineKeyboard()
		break
	case kbDataBack:
		var err error
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
		break
	case kbDataCancel:
		var err error
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
		err = s.processCancelCallback(ctx, bot, callbackQuery)
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
		markup, err = itemsListInlineKeyboard(ctx, s, messageId, "remove")
		if err != nil {
			return fmt.Errorf("failed to get markup for remove items list: %w", err)
		}
		break
	case kbDataEditItem:
		var err error
		markup, err = itemsListInlineKeyboard(ctx, s, messageId, "edit")
		if err != nil {
			return fmt.Errorf("failed to get markup for edit items list: %w", err)
		}
		break
	case kbDataCustomer:
		var err error
		markup, err = customerInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get markup for customer action: %w", err)
		}
		break
	case kbDataPayment:
		var err error
		markup, err = paymentInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get payment inlint markup: %w", err)
		}
		break
	case kbDataPaymentCard:
		var err error
		err = s.processAddPaymentCallback(ctx, callbackQuery, types.PaymentMethodCard2Card)
		if err != nil {
			return fmt.Errorf("failed to add card payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
		break
	case kbDataPaymentCash:
		var err error
		err = s.processAddPaymentCallback(ctx, callbackQuery, types.PaymentMethodCash)
		if err != nil {
			return fmt.Errorf("failed to add cash payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
		break
	case kbDataPaymentLink:
		var err error
		err = s.processAddPaymentCallback(ctx, callbackQuery, types.PaymentMethodAcquiring)
		if err != nil {
			return fmt.Errorf("failed to add link payment to order: %w", err)
		}
		markup, err = paymentAmountInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("faield to get payment amount inline markup: %w", err)
		}
		break
	case kbDataFullPayment:
		err := s.processPaymentCallback(ctx, bot, messageId, 0)
		if err != nil {
			return fmt.Errorf("failed to process full payment callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
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
		markup, err = paymentsListInlineKeyboard(ctx, s, messageId, "refund")
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
		order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order: %w", err)
		}
		err = s.processRefundCallback(ctx, bot, order, messageId, 0)
		if err != nil {
			return fmt.Errorf("failed to process refund callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case kbDataRemovePayment:
		var err error
		markup, err = paymentsListInlineKeyboard(ctx, s, messageId, "remove")
		if err != nil {
			return fmt.Errorf("failed to get payments list markup: %w", err)
		}
	case kbDataOrderActions:
		var err error
		markup, err = orderActionsInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order actions markup: %w", err)
		}
	case kbDataOrderDone:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.OrderStateDone)
		if err != nil {
			return fmt.Errorf("failed to process order done callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case kbDataOrderRestart:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.OrderStateForming)
		if err != nil {
			return fmt.Errorf("failed to process order restart callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case kbDataOrderDelete:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.OrderStateDeleted)
		if err != nil {
			return fmt.Errorf("failed to process order delete callback: %w", err)
		}
		markup = restoreDeletedOrderInlineKeyboard()
	case kbDataOrderRestore:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.OrderStateForming)
		if err != nil {
			return fmt.Errorf("failed to process order restore callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case kbDataOrderInProgress:
		err := s.processOrderStateCallback(ctx, bot, callbackQuery, types.OrderStateInProgress)
		if err != nil {
			return fmt.Errorf("failed to process order restore callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case kbDataOrderCollapse:
		err := s.processOrderDisplayModeCallback(ctx, bot, callbackQuery, types.DisplayModeCollapsed)
		if err != nil {
			return fmt.Errorf("failed to process order collapse callback: %w", err)
		}
		markup = expandOrderInlineKeyboard()
	case kbDataOrderExpand:
		err := s.processOrderDisplayModeCallback(ctx, bot, callbackQuery, types.DisplayModeFull)
		if err != nil {
			return fmt.Errorf("failed to process order expand callback: %w", err)
		}
		markup, err = startOrderInlineKeyboard(ctx, s, messageId)
		if err != nil {
			return fmt.Errorf("failed to get order inline kb: %w", err)
		}
	case kbDataDelivery:
		fallthrough
	case kbDataLingerieSet:
		err := s.processAddKnownItem(ctx, bot, callbackQuery, callbackData)
		if err != nil {
			return fmt.Errorf("failed to process add item callback: %w", err)
		}
		markup = cancelInlineKeyboard()
	case kbDataNotifyRead:
		shouldDelete = true
	default:
		args := strings.Split(callbackData, "_")
		entity := args[0]
		action := args[1]

		argsCount := len(args)

		if entity == "customer" {
			switch args[2] {
			case "email":
				markup = cancelInlineKeyboard()
				err := s.processCustomerEditEmail(ctx, bot, callbackQuery)
				if err != nil {
					return fmt.Errorf("failed to process customer edit name callback: %w", err)
				}
			case kbDataInstagram:
				markup = cancelInlineKeyboard()
				err := s.processCustomerEditInstagram(ctx, bot, callbackQuery)
				if err != nil {
					return fmt.Errorf("failed to process customer edit instagram callback: %w", err)
				}
			case "phone":
				markup = cancelInlineKeyboard()
				err := s.processCustomerEditPhone(ctx, bot, callbackQuery)
				if err != nil {
					return fmt.Errorf("faield to process customer edit phone callback: %w", err)
				}
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
				markup, err = startOrderInlineKeyboard(ctx, s, messageId)
				if err != nil {
					return fmt.Errorf("failed to get order inline kb: %w", err)
				}
				break
			case "refund":
				order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
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

	logrus.Info("switch case ended")

	var msg tg.Chattable
	if !shouldDelete {
		msg = tg.NewEditMessageReplyMarkup(chatId, int(messageId), markup)
	} else {
		msg = tg.NewDeleteMessage(chatId, int(messageId))
	}
	_, err := bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
