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

	ex := exchange.New(logger)

	go ex.Run()

	httpSrv := httpsrv.NewHTTPServer(logger, ex)

	go httpSrv.Run()

	<-quit

}
