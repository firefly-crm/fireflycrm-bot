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

func startOrderInlineKeyboard() tg.InlineKeyboardMarkup {
	itemsButton := tg.NewInlineKeyboardButtonData(kbItems, kbDataItems)
	customerButton := tg.NewInlineKeyboardButtonData(kbCustomer, kbDataCustomer)
	paymentButton := tg.NewInlineKeyboardButtonData(kbPayment, kbDataPayment)
	row1 := []tg.InlineKeyboardButton{
		itemsButton,
		customerButton,
		paymentButton,
	}
	actionsButton := tg.NewInlineKeyboardButtonData(kbOrderActions, kbDataOrderActions)
	row2 := []tg.InlineKeyboardButton{
		actionsButton,
	}

	return tg.NewInlineKeyboardMarkup(row1, row2)
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

	if order.State == types.StandBy {
		doneButton := tg.NewInlineKeyboardButtonData(kbOrderDone, kbDataOrderDone)
		rows = append(rows, []tg.InlineKeyboardButton{doneButton})
	}

	if order.State == types.Completed {
		restoreButton := tg.NewInlineKeyboardButtonData(kbOrderRestart, kbDataOrderRestart)
		rows = append(rows, []tg.InlineKeyboardButton{restoreButton})
	}

	deleteButton := tg.NewInlineKeyboardButtonData(kbOrderDelete, kbDataOrderDelete)
	rows = append(rows, []tg.InlineKeyboardButton{deleteButton})
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	rows = append(rows, []tg.InlineKeyboardButton{backButton})

	return tg.NewInlineKeyboardMarkup(rows...), nil
}

func orderItemsInlineKeyboard() tg.InlineKeyboardMarkup {
	addItemButton := tg.NewInlineKeyboardButtonData(kbAddItem, kbDataAddItem)
	editItemButton := tg.NewInlineKeyboardButtonData(kbEditItem, kbDataEditItem)
	removeItemButton := tg.NewInlineKeyboardButtonData(kbRemove, kbDataRemoveItem)
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	row1 := []tg.InlineKeyboardButton{addItemButton, removeItemButton}
	row2 := []tg.InlineKeyboardButton{editItemButton}
	row3 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row1, row2, row3)
}

func customerInlineKeyboard(ctx context.Context, s Service, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get order for markup: %w", err)
	}

	//nameButton := tg.NewInlineKeyboardButtonData(kbName, fmt.Sprintf("customer_edit_name_%d", order.Id))
	emailButton := tg.NewInlineKeyboardButtonData(kbCustomerEmail, fmt.Sprintf("customer_edit_email_%d", order.Id))
	//phoneButton := tg.NewInlineKeyboardButtonData(kbCustomerPhone, fmt.Sprintf("customer_edit_phone_%d", order.Id))
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)

	markups := [][]tg.InlineKeyboardButton{
		//{nameButton},
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
		if p.PaymentMethod == types.Acquiring && p.Payed {
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

func paymentInlineKeyboard() tg.InlineKeyboardMarkup {
	linkButton := tg.NewInlineKeyboardButtonData(kbPaymentLink, kbDataPaymentLink)
	cardButton := tg.NewInlineKeyboardButtonData(kbPaymentCard, kbDataPaymentCard)
	cashButton := tg.NewInlineKeyboardButtonData(kbPaymentCash, kbDataPaymentCash)
	row1 := []tg.InlineKeyboardButton{linkButton, cardButton, cashButton}
	deleteButton := tg.NewInlineKeyboardButtonData(kbRemove, kbDataRemovePayment)
	refundButton := tg.NewInlineKeyboardButtonData(kbRefundPayment, kbDataRefundPayment)
	row2 := []tg.InlineKeyboardButton{deleteButton, refundButton}
	cancelButton := tg.NewInlineKeyboardButtonData(kbCancel, kbDataCancel)
	row3 := []tg.InlineKeyboardButton{cancelButton}
	return tg.NewInlineKeyboardMarkup(row1, row2, row3)
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
