package main

import (
	"flag"
	"fmt"
	"github.com/DarthRamone/fireflycrm-bot/billmaker"
	"github.com/DarthRamone/fireflycrm-bot/infra"
	"github.com/DarthRamone/fireflycrm-bot/orderbook"
	"github.com/DarthRamone/fireflycrm-bot/service"
	"github.com/DarthRamone/fireflycrm-bot/storage"
	"github.com/DarthRamone/fireflycrm-bot/users"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

var token = flag.String("token", "", "Telegram bot token")

func main() {
	flag.Parse()

	tgToken := *token
	if tgToken == "" {
		tgToken = os.Getenv("TG_TOKEN")
		if tgToken == "" {
			panic("telegram token is unset. Use TG_TOKEN env, or --token cmd line arg")
		}
	}

	pgHost := os.Getenv("POSTGRES_HOST")
	if pgHost == "" {
		panic("pg host is unset; use POSTGRES_HOST env")
	}

	pgUser := os.Getenv("POSTGRES_USER")
	if pgUser == "" {
		panic("pg username is unset; use POSTGRES_USER env")
	}

	pgPassword := os.Getenv("POSTGRES_PASSWORD")
	if pgPassword == "" {
		panic("pg password is unset; use POSTGRES_PASSWORD env")
	}

	pgDBName := os.Getenv("POSTGRES_DB")
	if pgDBName == "" {
		panic("pg db is unset; user POSTGRES_DB env")
	}

	pgPort := "5432"
	envPort := os.Getenv("POSTGRES_PORT")
	if envPort != "" {
		pgPort = envPort
	}

	connString := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=disable", pgUser, pgPassword, pgDBName, pgPort, pgHost)

	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		logrus.Fatalln(err)
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

	err = g.Wait()
	if err != nil {
		panic(err)
	}
}
