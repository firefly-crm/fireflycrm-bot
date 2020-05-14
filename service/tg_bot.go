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

	wc := tg.NewWebhook("https://www.firefly.style/api/bot")
	_, err = bot.SetWebhook(wc)
	if err != nil {
		log.Fatalf("failed to set webhook: %v", err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	go func() {
		bot.Debug = false

		updates := bot.ListenForWebhook("/api/bot")

		for update := range updates {
			if ctx.Err() == context.Canceled {
				break
			}

			var err error

			info, err := bot.GetWebhookInfo()
			if err != nil {
				log.Fatal(err)
			}
			if info.LastErrorDate != 0 {
				log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
			}

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
