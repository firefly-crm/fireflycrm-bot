package service

import tg "github.com/go-telegram-bot-api/telegram-bot-api"

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
	removeItemButton := tg.NewInlineKeyboardButtonData(kbRemoveItem, kbDataRemoveItem)
	backButton := tg.NewInlineKeyboardButtonData(kbBack, kbDataBack)
	row1 := []tg.InlineKeyboardButton{addItemButton, removeItemButton}
	row2 := []tg.InlineKeyboardButton{backButton}
	return tg.NewInlineKeyboardMarkup(row1, row2)
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
