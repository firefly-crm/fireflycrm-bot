package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/orderbook"
	"github.com/DarthRamone/fireflycrm-bot/service"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/types"
	"github.com/DarthRamone/modulbank-go"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

var token = flag.String("token", "", "Telegram bot token")

func main() {
	flag.Parse()
	if *token == "" {
		panic("telegram bot token is unset")
	}

	db, err := sqlx.Connect("postgres", "user=admin password=1234 dbname=firefly port=32769 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	stor := storage.NewStorage(db)
	var mb modulbank.API

	ob := orderbook.MustNewOrderBook(stor, mb)

	serv := service.Service{
		OrderBook: ob,
	}

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Fatalf("failed to initialize bot: %w", err)
	}
	log.Printf("authorized on account %s", bot.Self.UserName)

	//bot.Debug = true
	updateConf := tgbotapi.NewUpdate(0)
	updateConf.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConf)
	if err != nil {
		log.Fatalf("failed to get updates channel: %w", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s (chat_id: %d)", update.Message.From.UserName, update.Message.Text, update.Message.Chat.ID)

		id, err := serv.OrderBook.CreateOrder(context.Background(), types.OrderOptions{
			Description:    "test",
			CustomerName:   "test name",
			CustomerEmail:  "test email",
			CustomerPhone:  "test phone",
			CustomerSocial: "http://insta.com",
		})
		if err != nil {
			fmt.Println(err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%d", id))
		msg.ReplyToMessageID = update.Message.MessageID
		button := tgbotapi.NewKeyboardButton("test")
		markup := tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{button})
		msg.ReplyMarkup = markup

		_, err = bot.Send(msg)
		if err != nil {
			log.Printf("failed to send message: %w", err)
		}
	}
}
