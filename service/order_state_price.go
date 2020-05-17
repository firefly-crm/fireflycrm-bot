package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func (s Service) setWaitingForPrice(ctx context.Context, bot *tg.BotAPI, order types.Order) error {
	err := s.OrderBook.UpdateOrderEditState(ctx, order.Id, types.EditStateWaitingItemPrice)
	if err != nil {
		return fmt.Errorf("failed to change order state: %w", err)
	}

	if !order.HintMessageId.Valid {
		return fmt.Errorf("hint message is nil")
	}

	hintMessageId := int(order.HintMessageId.Int64)

	hintMessage := tg.NewEditMessageText(int64(order.UserId), hintMessageId, replyEnterItemPrice)
	_, err = bot.Send(hintMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
