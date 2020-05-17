package service

import (
	"context"
)

type (
	Service struct {
		Publisher Publisher
	}

	Options struct {
		TelegramToken string
	}
)

func (s Service) Serve(ctx context.Context, opts Options) error {
	s.startListenTGUpdates(ctx, opts.TelegramToken)
	return nil
}
