package main

import (
	"context"
	"github.com/firefly-crm/common/infra"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/common/rabbit"
	"github.com/firefly-crm/common/rabbit/exchanges"
	"github.com/firefly-crm/fireflycrm-bot/config"
	"github.com/firefly-crm/fireflycrm-bot/service"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Crashf("service exited with error: %v", err)
		}
	}()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	serviceConfig := config.Config{}
	err = viper.Unmarshal(&serviceConfig)
	if err != nil {
		panic(err)
	}

	rabbitConfig := rabbit.Config{
		Endpoint: serviceConfig.Rabbit,
	}
	rabbitPrimary := rabbit.MustNew(rabbitConfig)
	exchange := exchanges.MustExchangeByID(exchanges.FireflyCRMTelegramUpdates)

	go func() {
		errPrimary := <-rabbitPrimary.Done()
		logger.Crashf("primary rabbit client error: %v", errPrimary)
	}()

	primaryPublisher := rabbitPrimary.MustNewExchange(exchange.Opts)

	serv := service.Service{
		Publisher: primaryPublisher,
	}

	ctx := infra.Context()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return serv.Serve(ctx, service.Options{TelegramToken: serviceConfig.TgToken})
	})

	err = g.Wait()
	if err != nil && err != context.Canceled {
		panic(err)
	}
}
