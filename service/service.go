package service

import (
	"context"
	"github.com/DarthRamone/fireflycrm-bot/billmaker"
	"github.com/DarthRamone/fireflycrm-bot/orderbook"
	"github.com/DarthRamone/fireflycrm-bot/users"
)

type (
	Service struct {
		OrderBook orderbook.OrderBook
		BillMaker billmaker.BillMaker
		Users     users.Users
	}

	ServiceOptions struct {
		TelegramToken string
	}
)

func (s Service) Serve(ctx context.Context, opts ServiceOptions) {
	s.startListenTGUpdates(ctx, opts.TelegramToken)
}
