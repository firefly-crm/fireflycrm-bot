package service

import (
	"context"
	"fmt"
	. "github.com/firefly-crm/common/bot"
	"github.com/firefly-crm/common/logger"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/common/rabbit/routes"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func (s Service) processCallback(ctx context.Context, update tg.Update) error {
	callbackQuery := update.CallbackQuery
	userId := uint64(callbackQuery.Message.Chat.ID)
	messageId := uint64(callbackQuery.Message.MessageID)

	callbackData := callbackQuery.Data

	log := logger.FromContext(ctx).
		WithField("user_id", userId).
		WithField("callback", callbackData).
		WithField("message_id", messageId)

	ctx = logger.ToContext(ctx, log)

	log.Infof("processing callback")

	var callbackType tp.CallbackType
	var entityId uint64

	switch callbackData {
	case KbDataItems:
		callbackType = tp.CallbackType_ITEMS
	case KbDataBack:
		callbackType = tp.CallbackType_BACK
	case KbDataCancel:
		callbackType = tp.CallbackType_CANCEL
	case KbDataAddItem:
		callbackType = tp.CallbackType_RECEIPT_ITEMS_ADD
	case KbDataRemoveItem:
		callbackType = tp.CallbackType_RECEIPT_ITEMS_REMOVE
	case KbDataEditItem:
		callbackType = tp.CallbackType_RECEIPT_ITEMS_EDIT
	case KbDataCustomer:
		callbackType = tp.CallbackType_CUSTOMER
	case KbDataPayment:
		callbackType = tp.CallbackType_PAYMENTS
	case KbDataPaymentCard:
		callbackType = tp.CallbackType_ADD_PAYMENT_TRANSFER
	case KbDataPaymentCash:
		callbackType = tp.CallbackType_ADD_PAYMENT_CASH
	case KbDataPaymentLink:
		callbackType = tp.CallbackType_ADD_PAYMENT_LINK
	case KbDataFullPayment:
		callbackType = tp.CallbackType_PAYMENT_AMOUNT_FULL
	case KbDataPartialPayment:
		callbackType = tp.CallbackType_PAYMENT_AMOUNT_PARTIAL
	case KbDataRefundPayment:
		callbackType = tp.CallbackType_PAYMENTS_REFUND
	case KbDataPartialRefund:
		callbackType = tp.CallbackType_PAYMENT_REFUND_PARTIAL
	case KbDataFullRefund:
		callbackType = tp.CallbackType_PAYMENT_REFUND_FULL
	case KbDataRemovePayment:
		callbackType = tp.CallbackType_PAYMENTS_REMOVE
	case KbDataOrderActions:
		callbackType = tp.CallbackType_ORDER_ACTIONS
	case KbDataOrderDone:
		callbackType = tp.CallbackType_ORDER_STATE_DONE
	case KbDataOrderRestart:
		callbackType = tp.CallbackType_ORDER_RESTART
	case KbDataOrderDelete:
		callbackType = tp.CallbackType_ORDER_DELETE
	case KbDataOrderRestore:
		callbackType = tp.CallbackType_ORDER_RESTORE
	case KbDataOrderInProgress:
		callbackType = tp.CallbackType_ORDER_STATE_IN_PROGRESS
	case KbDataOrderCollapse:
		callbackType = tp.CallbackType_ORDER_COLLAPSE
	case KbDataOrderExpand:
		callbackType = tp.CallbackType_ORDER_EXPAND
	case KbDataDelivery:
		callbackType = tp.CallbackType_CUSTOM_ITEM_DELIVERY
	case KbDataLingerieSet:
		callbackType = tp.CallbackType_CUSTOM_ITEM_LINGERIE_SET
	case KbDataNotifyRead:
		callbackType = tp.CallbackType_NOTIFY_READ
	case KbDataOrderEdit:
		callbackType = tp.CallbackType_ORDER_EDIT
	default:
		args := strings.Split(callbackData, "_")
		entity := args[0]
		action := args[1]

		argsCount := len(args)

		if entity == "customer" {
			switch args[2] {
			case "email":
				callbackType = tp.CallbackType_CUSTOMER_EDIT_EMAIL
			case "instagram":
				callbackType = tp.CallbackType_CUSTOMER_EDIT_INSTAGRAM
			case "phone":
				callbackType = tp.CallbackType_CUSTOMER_EDIT_PHONE
			case "description":
				callbackType = tp.CallbackType_CUSTOMER_EDIT_DESCRIPTION
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
				callbackType = tp.CallbackType_PAYMENT_REMOVE
				entityId = id
			case "refund":
				callbackType = tp.CallbackType_PAYMENT_REFUND
				entityId = id
			}
		}

		if entity == "order" {
			var err error
			entityId, err = strconv.ParseUint(args[3], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse id: %w", err)
			}

			if action == "edit" {
				switch args[2] {
				case "date":
					callbackType = tp.CallbackType_ORDER_EDIT_DUE_DATE
				case "description":
					callbackType = tp.CallbackType_ORDER_EDIT_DESCRIPTION
				}
			}
		}

		if entity == "item" {
			if argsCount == 3 {
				var err error
				entityId, err = strconv.ParseUint(args[2], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse id: %w", err)
				}

				if action == "edit" {
					callbackType = tp.CallbackType_RECEIPT_ITEM_EDIT
				}

				if action == "remove" {
					callbackType = tp.CallbackType_RECEIPT_ITEM_REMOVE
				}
			}

			if argsCount == 4 {
				var err error
				entityId, err = strconv.ParseUint(args[3], 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse id: %w", err)
				}

				switch args[2] {
				case "qty":
					callbackType = tp.CallbackType_RECEIPT_ITEM_EDIT_QTY
				case "name":
					callbackType = tp.CallbackType_RECEIPT_ITEM_EDIT_NAME
				case "price":
					callbackType = tp.CallbackType_RECEIPT_ITEM_EDIT_PRICE
				}
			}
		}
	}

	callbackEvent := &tp.CallbackEvent{
		UserId:    userId,
		MessageId: messageId,
		EntityId:  entityId,
		Event:     callbackType,
	}

	rt, err := routes.QueueByID(routes.TelegramCallbackUpdate)
	if err != nil {
		return fmt.Errorf("failed to get route: %w", err)
	}

	err = s.Publisher.Publish(ctx, rt.Route(), callbackEvent)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
