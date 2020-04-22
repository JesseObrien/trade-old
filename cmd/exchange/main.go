package main

import (
	"os"
	"os/signal"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jesseobrien/trade/internal/exchange"
	"github.com/jesseobrien/trade/internal/httpsrv"
)

func main() {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	defer signal.Stop(quit)

	logger := log.Logger{
		Handler: cli.New(os.Stdout),
	}

	logger.Info("📈 Welcome to Trade 📈")

	exchange := exchange.New(logger)

	go exchange.Run()

	httpSrv := httpsrv.NewHTTPServer(logger)

	go httpSrv.Run()

	<-quit

}
