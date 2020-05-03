package main

import (
	"context"
	"flag"
	"github.com/DarthRamone/fireflycrm-bot/billmaker"
	"github.com/DarthRamone/fireflycrm-bot/orderbook"
	"github.com/DarthRamone/fireflycrm-bot/service"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/users"
	"github.com/DarthRamone/modulbank-go"
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
	bm := billmaker.NewBillMaker()
	u := users.NewUsers(stor)

	serv := service.Service{
		OrderBook: ob,
		BillMaker: bm,
		Users:     u,
	}

	ctx := context.Background()

	serv.Serve(ctx, service.ServiceOptions{TelegramToken: *token})
}
