package service

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func merchantStandByKeyboardMarkup() tg.ReplyKeyboardMarkup {
	createOrderButton := tg.NewKeyboardButton(kbCreateOrder)
	return tg.NewReplyKeyboard([]tg.KeyboardButton{createOrderButton})
}

func startOrderInlineKeyboard() tg.InlineKeyboardMarkup {
	itemsButton := tg.NewInlineKeyboardButtonData(kbItems, kbDataItems)
	customerButton := tg.NewInlineKeyboardButtonData(kbCustomer, kbDataCustomer)
	return tg.NewInlineKeyboardMarkup([]tg.InlineKeyboardButton{itemsButton, customerButton})
}

func orderItemsInlineKeyboard() tg.InlineKeyboardMarkup {
	addItemButton := tg.NewInlineKeyboardButtonData(kbAddItem, kbDataAddItem)
	editItemButton := tg.NewInlineKeyboardButtonData(kbEditItem, kbDataEditItem)
	removeItemButton := tg.NewInlineKeyboardButtonData(kbRemoveItem, kbDataRemoveItem)
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	row1 := []tg.InlineKeyboardButton{addItemButton, editItemButton, removeItemButton}
	row2 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row1, row2)
}

func itemsListInlineKeyboard(ctx context.Context, s Service, messageId uint64) (tg.InlineKeyboardMarkup, error) {
	var markup tg.InlineKeyboardMarkup

	order, err := s.OrderBook.GetOrderByMessageId(ctx, messageId)
	if err != nil {
		return markup, fmt.Errorf("failed to get items markup: %w", err)
	}

	markups := make([][]tg.InlineKeyboardButton, 0)
	for _, i := range order.ReceiptItems {
		button := tg.NewInlineKeyboardButtonData(i.Name, fmt.Sprintf("item_edit_%d", i.Id))
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

func addItemInlineKeyboard() tg.InlineKeyboardMarkup {
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	row1 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row1)
}

func cancelInlineKeyboard() tg.InlineKeyboardMarkup {
	cancelButton := tg.NewInlineKeyboardButtonData(kbCancel, kbDataCancel)
	row1 := []tg.InlineKeyboardButton{cancelButton}
	return tg.NewInlineKeyboardMarkup(row1)
}
