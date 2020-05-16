package service

import (
	"context"
	tg "github.com/DarthRamone/telegram-bot-api"
	"github.com/sirupsen/logrus"
)

func (s Service) startListenTGUpdates(ctx context.Context, token string) *tg.BotAPI {
	bot, err := tg.NewBotAPI(token)
	if err != nil {
		logrus.Errorf("failed to initialize bot: %v", err)
	}
	logrus.Infof("authorized on account %s", bot.Self.UserName)

	wc := tg.NewWebhook("https://www.firefly.style/api/bot")
	_, err = bot.SetWebhook(wc)
	if err != nil {
		logrus.Errorf("failed to set webhook: %v", err)
	}

	//_, err = bot.RemoveWebhook()
	//if err != nil {
	//	logrus.Errorf("failed to remove webhook: %v", err)
	//}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		logrus.Error(err)
	}
	if info.LastErrorDate != 0 {
		logrus.Warnf("Telegram callback failed: %s; date: %", info.LastErrorMessage)
	}

	go func() {
		bot.Debug = false

		updates := bot.ListenForWebhook("/api/bot")
		u := tg.NewUpdate(0)
		u.Timeout = 60

		//updates, _ := bot.GetUpdatesChan(u)
		for update := range updates {
			if ctx.Err() == context.Canceled {
				break
			}

			logrus.Info("update received")

			var err error
			info, err := bot.GetWebhookInfo()
			if err != nil {
				logrus.Fatal(err)
			}
			if info.LastErrorDate != 0 {
				logrus.Warnf("Telegram callback failed: %s", info.LastErrorMessage)
			}

			if update.CallbackQuery != nil {
				logrus.Info("processing callback")
				err = s.processCallback(ctx, bot, update)
			} else {
				if update.Message == nil {
					continue
				}

				logrus.Info("processing command")
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
