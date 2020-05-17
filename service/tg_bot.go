package service

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/common/logger"
	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
)

const (
	WEBHOOK_URL  = "https://www.firefly.style/api/bot"
	WEBHOOK_PATH = "/api/bot"
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

	wc := tg.NewWebhook(WEBHOOK_URL)
	_, err = bot.SetWebhook(wc)
	if err != nil {
		log.Errorf("failed to set webhook: %v", err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Errorf("failed to get webhook info: %v", err)
	}
	if info.LastErrorDate != 0 {
		log.Warnf("telegram webhook last error: %s", info.LastErrorMessage)
	}

	go func() {
		bot.Debug = false

		updates := bot.ListenForWebhook(WEBHOOK_PATH)

		for update := range updates {
			if ctx.Err() == context.Canceled {
				break
			}
			log.Infof("update received")

			var err error
			info, err := bot.GetWebhookInfo()
			if err != nil {
				log.Errorf(err)
			}
			if info.LastErrorDate != 0 {
				log.Warnf("telegram webhook last error: %s", info.LastErrorMessage)
			}

			if update.CallbackQuery != nil {
				err = s.processCallback(ctx, bot, update)
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
