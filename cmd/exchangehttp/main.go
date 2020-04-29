package main

import (
	"os"
	"os/signal"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jesseobrien/trade/internal/httpsrv"
	"github.com/jesseobrien/trade/internal/nats"
)

func main() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	logger := log.Logger{
		Handler: cli.New(os.Stdout),
	}

	logger.Info("connecting to nats server")
	natsConn, err := nats.NewJSONConnection()
	defer natsConn.Close()

	if err != nil {
		logger.Errorf("could not connect to nats %v", err)
		panic(err)
	}

	logger.Info("nats connected")

	httpSrv := httpsrv.NewHTTPServer(logger, natsConn)

	go httpSrv.Run()

	<-quit

	logger.Info("Received interrupt, shutting down")
}
