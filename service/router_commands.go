package service

import (
	"context"
	"fmt"
	. "github.com/firefly-crm/common/bot"
	"github.com/firefly-crm/common/logger"
	tp "github.com/firefly-crm/common/messages/telegram"
	"github.com/firefly-crm/common/rabbit/routes"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func (s Service) processCommand(ctx context.Context, bot *tg.BotAPI, update tg.Update) error {
	var err error
	var cmd = update.Message.Text

	log := logger.
		FromContext(ctx).
		WithField("user_id", update.Message.From.ID).
		WithField("command", cmd)

	log.Infof("processing command")

	ctx = logger.ToContext(ctx, log)

	var commandType tp.CommandType
	args := make([]string, 0)

	if cmd == "/start" {
		commandType = tp.CommandType_START
	} else if cmd == KbCreateOrder {
		commandType = tp.CommandType_CREATE_ORDER
	} else if strings.HasPrefix(cmd, "/registerAsMerchant") {
		commandType = tp.CommandType_REGISTER_AS_MERCHANT
		args = strings.Split(cmd, " ")[1:]
	} else if cmd == KbActiveOrders {
	} else {
		promptEvent := &tp.PromptEvent{
			UserId:      uint64(update.Message.From.ID),
			MessageId:   uint64(update.Message.MessageID),
			UserMessage: cmd,
		}

		rt, err := routes.QueueByID(routes.TelegramPromptUpdate)
		if err != nil {
			return fmt.Errorf("failed to get route: %w", err)
		}

		err = s.Publisher.Publish(ctx, rt.Route(), promptEvent)
		if err != nil {
			return fmt.Errorf("failed to publish message: %w", err)
		}

		return nil
	}

	promptEvent := &tp.CommandEvent{
		UserId:    uint64(update.Message.From.ID),
		MessageId: uint64(update.Message.MessageID),
		Arguments: args,
		Command:   commandType,
	}

	rt, err := routes.QueueByID(routes.TelegramCommandUpdate)
	if err != nil {
		return fmt.Errorf("failed to get route: %w", err)
	}

	err = s.Publisher.Publish(ctx, rt.Route(), promptEvent)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
