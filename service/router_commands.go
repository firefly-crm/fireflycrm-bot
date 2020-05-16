package service

import (
	"context"
	"fmt"
	tg "github.com/DarthRamone/telegram-bot-api"
	"strings"
)

func (s Service) processCommand(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	var err error
	var cmd = update.Message.Text

	if cmd == "/start" {
		err = s.createUser(ctx, bot, update)
	} else if cmd == kbCreateOrder {
		err = s.createOrder(ctx, bot, update)
	} else if strings.HasPrefix(cmd, "/registerAsMerchant") {
		err = s.registerMerchant(ctx, bot, update)
	} else if cmd == kbActiveOrders {

	} else {
		err = s.processPrompt(ctx, bot, update)
	}

	if err != nil {
		return fmt.Errorf("failed process message: %w", err)
	}

	return nil
}
