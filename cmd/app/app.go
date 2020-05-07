package main

import (
	"flag"
	"github.com/DarthRamone/fireflycrm-bot/billmaker"
	"github.com/DarthRamone/fireflycrm-bot/infra"
	"github.com/DarthRamone/fireflycrm-bot/orderbook"
	"github.com/DarthRamone/fireflycrm-bot/service"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/users"
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
	ob := orderbook.MustNewOrderBook(stor)
	bm := billmaker.NewBillMaker()
	u := users.NewUsers(stor)

	serv := service.Service{
		OrderBook: ob,
		BillMaker: bm,
		Users:     u,
		Storage:   stor,
	}

	ctx := infra.Context()
	serv.Serve(ctx, service.Options{TelegramToken: *token})
}
