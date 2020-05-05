package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"sort"
)

func merchantStandByKeyboardMarkup() tg.ReplyKeyboardMarkup {
	createOrderButton := tg.NewKeyboardButton(kbCreateOrder)
	getOrdersButton := tg.NewKeyboardButton(kbActiveOrders)
	rows := [][]tg.KeyboardButton{
		{createOrderButton},
		{getOrdersButton},
	}
	return tg.NewReplyKeyboard(rows...)
}

func expandOrderInlineKeyboard() tg.InlineKeyboardMarkup {
	expandButton := tg.NewInlineKeyboardButtonData(kbOrderExpand, kbDataOrderExpand)
	return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(expandButton))
}

func startOrderInlineKeyboard(ctx context.Context, s Service, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	customerButton := tg.NewInlineKeyboardButtonData(kbCustomer, kbDataCustomer)
	paymentButton := tg.NewInlineKeyboardButtonData(kbPayment, kbDataPayment)
	actionsButton := tg.NewInlineKeyboardButtonData(kbOrderActions, kbDataOrderActions)
	itemsButton := tg.NewInlineKeyboardButtonData(kbItems, kbDataItems)

	if order.OrderState != types.OrderStateDone {
		row1 := []tg.InlineKeyboardButton{itemsButton, customerButton, paymentButton}
		row2 := []tg.InlineKeyboardButton{actionsButton}
		markup = tg.NewInlineKeyboardMarkup(row1, row2)
	} else {
		row1 := []tg.InlineKeyboardButton{customerButton, paymentButton, actionsButton}
		markup = tg.NewInlineKeyboardMarkup(row1)
	}

	return markup, nil
}

func restoreDeletedOrderInlineKeyboard() tg.InlineKeyboardMarkup {
	restoreButton := tg.NewInlineKeyboardButtonData(kbOrderRestore, kbDataOrderRestore)
	return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(restoreButton))
}

func orderActionsInlineKeyboard(ctx context.Context, s Service, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup
	var rows [][]tg.InlineKeyboardButton

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	doneButton := tg.NewInlineKeyboardButtonData(kbOrderDone, kbDataOrderDone)
	inProgressButton := tg.NewInlineKeyboardButtonData(kbOrderInProgress, kbDataOrderInProgress)
	collapseButton := tg.NewInlineKeyboardButtonData(kbOrderCollapse, kbDataOrderCollapse)
	restartButton := tg.NewInlineKeyboardButtonData(kbOrderRestart, kbDataOrderRestart)
	deleteButton := tg.NewInlineKeyboardButtonData(kbOrderDelete, kbDataOrderDelete)
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)

	if order.OrderState == types.OrderStateForming && order.Amount > 0 {
		row := []tg.InlineKeyboardButton{inProgressButton}
		if order.PayedAmount == order.Amount {
			row = append(row, doneButton)
		}
		rows = append(rows, row)
	}

	if order.OrderState == types.OrderStateInProgress {
		if order.PayedAmount == order.Amount {
			rows = append(rows, []tg.InlineKeyboardButton{doneButton})
		}
	}

	rows = append(rows, []tg.InlineKeyboardButton{collapseButton})

	if order.OrderState == types.OrderStateDone {
		rows = append(rows, []tg.InlineKeyboardButton{restartButton})
	}

	rows = append(rows, []tg.InlineKeyboardButton{deleteButton})
	rows = append(rows, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(rows...), nil
}

func orderItemsInlineKeyboard() tg.InlineKeyboardMarkup {
	addItemButton := tg.NewInlineKeyboardButtonData(kbAddItemPictogram, kbDataAddItem)
	editItemButton := tg.NewInlineKeyboardButtonData(kbEditItemPictogram, kbDataEditItem)
	removeItemButton := tg.NewInlineKeyboardButtonData(kbRemovePictogram, kbDataRemoveItem)
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	deliveryButton := tg.NewInlineKeyboardButtonData(kbDeliveryPictogram, kbDataDelivery)
	lingerieButton := tg.NewInlineKeyboardButtonData(kbLingerieSetPictogram, kbDataLingerieSet)
	row1 := []tg.InlineKeyboardButton{addItemButton, deliveryButton, lingerieButton, editItemButton, removeItemButton}
	row2 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row1, row2)
}

func customerInlineKeyboard(ctx context.Context, s Service, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	//nameButton := tg.NewInlineKeyboardButtonData(kbName, fmt.Sprintf("customer_edit_name_%d", order.Id))
	emailButton := tg.NewInlineKeyboardButtonData(kbCustomerEmail, fmt.Sprintf("customer_edit_email_%d", order.Id))
	instaButton := tg.NewInlineKeyboardButtonData(kbCustomerInstagram, fmt.Sprintf("customer_edit_instagram_%d", order.Id))
	//phoneButton := tg.NewInlineKeyboardButtonData(kbCustomerPhone, fmt.Sprintf("customer_edit_phone_%d", order.Id))
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)

	markups := [][]tg.InlineKeyboardButton{
		//{nameButton},
		{instaButton},
		{emailButton},
		//{phoneButton},
		{backButton},
	}

	markup = tg.NewInlineKeyboardMarkup(markups...)
	return markup, nil
}

func itemsListInlineKeyboard(ctx context.Context, s Service, messageId uint64, action string) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get items markup: %w", err)
	}

	markups := make([][]tg.InlineKeyboardButton, 0)
	for _, i := range order.ReceiptItems {
		button := tg.NewInlineKeyboardButtonData(i.Name, fmt.Sprintf("item_%s_%d", action, i.Id))
		markups = append(markups, []tg.InlineKeyboardButton{button})
	}
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	markups = append(markups, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(markups...), nil
}

func paymentsListInlineKeyboard(ctx context.Context, s Service, messageId uint64, action string) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get payments markup: %w", err)
	}

	sort.Sort(types.PaymentsByCreatedAt(order.Payments))

	markups := make([][]tg.InlineKeyboardButton, 0)
	for i, p := range order.Payments {
		if p.PaymentMethod == types.PaymentMethodAcquiring && p.Payed {
			continue
		}

		if action == "refund" && !p.Payed {
			continue
		}

		if action == "refund" && p.RefundAmount == p.Amount {
			continue
		}

		name := fmt.Sprintf("Платеж #%d", i+1)
		button := tg.NewInlineKeyboardButtonData(name, fmt.Sprintf("payment_%s_%d", action, p.Id))
		markups = append(markups, []tg.InlineKeyboardButton{button})
	}
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	markups = append(markups, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(markups...), nil
}

func editItemActionsInlineKeyboard(itemId uint64) tg.InlineKeyboardMarkup {
	nameButton := tg.NewInlineKeyboardButtonData(kbName, fmt.Sprintf("item_edit_name_%d", itemId))
	qtyButton := tg.NewInlineKeyboardButtonData(kbQty, fmt.Sprintf("item_edit_qty_%d", itemId))
	priceButton := tg.NewInlineKeyboardButtonData(kbPrice, fmt.Sprintf("item_edit_price_%d", itemId))
	row := []tg.InlineKeyboardButton{nameButton, qtyButton, priceButton}
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	row1 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row, row1)
}

func cancelInlineKeyboard() tg.InlineKeyboardMarkup {
	cancelButton := tg.NewInlineKeyboardButtonData(kbCancel, kbDataCancel)
	row1 := []tg.InlineKeyboardButton{cancelButton}
	return tg.NewInlineKeyboardMarkup(row1)
}

func paymentInlineKeyboard(ctx context.Context, s Service, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get payments markup: %w", err)
	}

	linkButton := tg.NewInlineKeyboardButtonData(kbPaymentLink, kbDataPaymentLink)
	cardButton := tg.NewInlineKeyboardButtonData(kbPaymentCard, kbDataPaymentCard)
	cashButton := tg.NewInlineKeyboardButtonData(kbPaymentCash, kbDataPaymentCash)
	deleteButton := tg.NewInlineKeyboardButtonData(kbRemove, kbDataRemovePayment)
	refundButton := tg.NewInlineKeyboardButtonData(kbRefundPayment, kbDataRefundPayment)
	cancelButton := tg.NewInlineKeyboardButtonData(kbCancel, kbDataCancel)

	if order.Amount == 0 {
		return tg.NewInlineKeyboardMarkup(tg.NewInlineKeyboardRow(cancelButton)), nil
	}

	if order.PayedAmount >= order.Amount {
		var row1 []tg.InlineKeyboardButton
		if order.OrderState == types.OrderStateDone {
			row1 = []tg.InlineKeyboardButton{refundButton}
		} else {
			row1 = []tg.InlineKeyboardButton{deleteButton, refundButton}
		}
		row2 := []tg.InlineKeyboardButton{cancelButton}
		markup = tg.NewInlineKeyboardMarkup(row1, row2)
	} else {
		row1 := []tg.InlineKeyboardButton{linkButton, cardButton, cashButton}
		row2 := []tg.InlineKeyboardButton{deleteButton, refundButton}
		row3 := []tg.InlineKeyboardButton{cancelButton}
		markup = tg.NewInlineKeyboardMarkup(row1, row2, row3)
	}

	return markup, nil
}

func paymentAmountInlineKeyboard() tg.InlineKeyboardMarkup {
	fullButton := tg.NewInlineKeyboardButtonData(kbFullPayment, kbDataFullPayment)
	partialButton := tg.NewInlineKeyboardButtonData(kbPartialPayment, kbDataPartialPayment)
	cancelButton := tg.NewInlineKeyboardButtonData(kbCancel, kbDataCancel)

	markups := [][]tg.InlineKeyboardButton{
		{fullButton},
		{partialButton},
		{cancelButton},
	}
	return tg.NewInlineKeyboardMarkup(markups...)
}

func refundAmountInlineKeyboard() tg.InlineKeyboardMarkup {
	fullButton := tg.NewInlineKeyboardButtonData(kbFullRefund, kbDataFullRefund)
	partialButton := tg.NewInlineKeyboardButtonData(kbPartialRefund, kbDataPartialRefund)
	cancelButton := tg.NewInlineKeyboardButtonData(kbCancel, kbDataCancel)

	markups := [][]tg.InlineKeyboardButton{
		{fullButton},
		{partialButton},
		{cancelButton},
	}
	return tg.NewInlineKeyboardMarkup(markups...)
}
