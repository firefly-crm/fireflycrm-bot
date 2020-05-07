package service

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/billmaker"
	"github.com/DarthRamone/fireflycrm-bot/orderbook"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/users"
)

type (
	Service struct {
		OrderBook orderbook.OrderBook
		BillMaker billmaker.BillMaker
		Users     users.Users
		Storage   storage.Storage
	}

	Options struct {
		TelegramToken string
	}
)

func (s Service) Serve(ctx context.Context, opts Options) {
	bot := s.startListenTGUpdates(ctx, opts.TelegramToken)
	s.startPaymentsWatcher(ctx, bot)
}
