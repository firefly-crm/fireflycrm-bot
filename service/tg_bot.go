package service

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
)

func (s Service) startListenTGUpdates(ctx context.Context, token string) {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to initialize bot: %w", err)
	}
	log.Printf("authorized on account %s", bot.Self.UserName)

	//bot.Debug = true
	updateConf := tg.NewUpdate(0)
	updateConf.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConf)
	if err != nil {
		log.Fatalf("failed to get updates channel: %w", err)
	}

	for update := range updates {
		var err error

		if update.CallbackQuery != nil {
			log.Println("callback is not null")
			err = s.processCallback(ctx, bot, update)
		} else {
			if update.Message == nil {
				continue
			}

			err = s.processCommand(ctx, bot, update)
		}

		if err != nil {
			//TODO: Restore state
			logrus.Error(ctx, "failed to process message: %v", err.Error())
		}
	}
}
