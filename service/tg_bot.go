package service

import (
	"context"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"log"
)

func (s Service) startListenTGUpdates(ctx context.Context, token string) *tg.BotAPI {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		log.Fatalf("failed to initialize bot: %v", err)
	}
	log.Printf("authorized on account %s", bot.Self.UserName)

	go func() {
		//bot.Debug = true
		updateConf := tg.NewUpdate(0)
		updateConf.Timeout = 60

		updates, err := bot.GetUpdatesChan(updateConf)
		if err != nil {
			log.Fatalf("failed to get updates channel: %w", err)
		}

		for update := range updates {
			if ctx.Err() == context.Canceled {
				break
			}

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
				logrus.Errorf("failed to process message: %v", err.Error())
			}
		}
	}()

	return bot
}
