package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jesseobrien/trade/internal/httpsrv"
	"github.com/jesseobrien/trade/internal/nats"
)

var natsURL = "nats://localhost:4222"

func init() {
	flag.StringVar(&natsURL, "nats", "nats://localhost:4222", "-nats nats://localhost:4222")
	flag.Parse()
}

func main() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	logger := log.Logger{
		Handler: cli.New(os.Stdout),
	}

	logger.Info("✉️ connecting to nats server")
	natsConn, err := nats.NewJSONConnection(natsURL)

	if err != nil {
		logger.Errorf("could not connect to nats %v", err)
		panic(err)
	}

	defer natsConn.Close()

	logger.Info("✉️ nats connected")

	httpSrv := httpsrv.NewHTTPServer(logger, natsConn)

	go httpSrv.Run()

	<-quit

	logger.Info("Received interrupt, shutting down")
}
