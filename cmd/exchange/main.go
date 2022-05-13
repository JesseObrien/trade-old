package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jesseobrien/trade/internal/exchange"
	"github.com/jesseobrien/trade/internal/nats"
	"github.com/shopspring/decimal"
)

var natsURL = "nats://localhost:4222"

func init() {
	flag.StringVar(&natsURL, "nats", "nats://localhost:4222", "-nats nats://localhost:4222")
	flag.Parse()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)

	logger := log.Logger{
		Handler: cli.New(os.Stdout),
	}

	natsConn, err := nats.NewJSONConnection(natsURL)

	if err != nil {
		logger.Errorf("NATS connection error %v", err)
		panic(err)
	}

	logger.Info("📈 Welcome to Trade 📈")

	ex := exchange.New(logger, natsConn)

	symbols, err := os.Open("symbols.txt")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(symbols)

	for scanner.Scan() {
		ob := orderbook.NewOrderBook(logger, exchange.Symbol(scanner.Text()))
		//min + rand.Float64() * (max - min)
		randPrice := 1.0 + rand.Float64()*(100.0-1.0)
		ex.IPO(ob, decimal.NewFromFloat(randPrice), 1000)
	}

	go ex.Run()

	<-quit
}
