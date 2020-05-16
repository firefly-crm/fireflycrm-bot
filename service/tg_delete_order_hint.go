package service

import (
	"context"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/types"
	tg "github.com/DarthRamone/telegram-bot-api"
)

func (s Service) deleteHint(ctx context.Context, bot *tg.BotAPI, order types.Order) error {
	if !order.HintMessageId.Valid {
		return fmt.Errorf("hint id is nil")
	}

	deleteMessage := tg.NewDeleteMessage(int64(order.UserId), int(order.HintMessageId.Int64))
	_, err := bot.Send(deleteMessage)
	if err != nil {
		return fmt.Errorf("failed to delete hind: %w", err)
	}

	return nil
}
