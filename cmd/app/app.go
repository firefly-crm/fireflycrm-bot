package main

import (
	"flag"
	"fmt"
	"github.com/firefly-crm/common/infra"
	"github.com/firefly-crm/common/logger"
	"github.com/firefly-crm/common/rabbit"
	"github.com/firefly-crm/common/rabbit/exchanges"
	"github.com/firefly-crm/fireflycrm-bot/service"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
)

var token = flag.String("token", "", "Telegram bot token")

func main() {
	defer func() {
		if err := recover(); err != nil {
			logger.Crashf("service exited with error: %v", err)
		}
	}()

	flag.Parse()

	tgToken := *token
	if tgToken == "" {
		tgToken = os.Getenv("TG_TOKEN")
		if tgToken == "" {
			panic("telegram token is unset. Use TG_TOKEN env, or --token cmd line arg")
		}
	}

	rabbitUsername := os.Getenv("RMQ_USERNAME")
	if rabbitUsername == "" {
		log.Fatalf("rabbit username is empty")
	}

	rabbitPassword := os.Getenv("RMQ_PASSWORD")
	if rabbitPassword == "" {
		log.Fatal("rabbit password is empty")
	}

	rabbitHost := os.Getenv("RMQ_HOST")
	if rabbitHost == "" {
		log.Fatalf("rabbit host is empty")
	}

	rabbitPort := os.Getenv("RMQ_PORT")
	if rabbitPort == "" {
		log.Fatalf("rabbit port is empty")
	}

	rabbitConnString := fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitUsername, rabbitPassword, rabbitHost, rabbitPort)

	rabbitConfig := rabbit.Config{
		Endpoint: rabbitConnString,
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("FIREFLY LINGERIE"))
		if err != nil {
			fmt.Printf("handle err: %w", err)
		}
	})
	http.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello world"))
		if err != nil {
			fmt.Printf("handle err: %w", err)
		}
	})

	g.Go(func() error {
		return serv.Serve(ctx, service.Options{TelegramToken: *token})
	})
	g.Go(func() error {
		return http.ListenAndServe(":80", nil)
	})

	err := g.Wait()
	if err != nil {
		panic(err)
	}
}
