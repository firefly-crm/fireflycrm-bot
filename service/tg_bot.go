package service

import (
	"context"
	"github.com/firefly-crm/common/logger"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

func (s Service) startListenTGUpdates(ctx context.Context, token string) *tg.BotAPI {
	log := logger.FromContext(ctx)

	var c *http.Client
	transport, ok := http.DefaultTransport.(*http.Transport)
	if ok {
		transport.DisableKeepAlives = true
		var rt http.RoundTripper = transport
		c = &http.Client{Transport: rt}
	} else {
		c = http.DefaultClient
	}

	bot, err := tg.NewBotAPIWithClient(token, c)
	if err != nil {
		log.Errorf("failed to initialize bot: %v", err)
	}
	log.Infof("authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Errorf("failed to get updates chan: %w", err)
		return nil
	}

	go func() {
		bot.Debug = false

		for update := range updates {
			if ctx.Err() == context.Canceled {
				break
			}
			log.Infof("update received")

			if update.CallbackQuery != nil {
				err = s.processCallback(ctx, update)
			} else {
				if update.Message == nil {
					continue
				}
				err = s.processCommand(ctx, bot, update)
			}

			if err != nil {
				log.Errorf("failed to process message: %v", err.Error())
			}
		}
	}()

	return bot
}
