package main

import (
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/jesseobrien/trade/internal/httpsrv"
	"github.com/jesseobrien/trade/internal/service"
)

func main() {

	logger := log.Logger{
		Handler: cli.New(os.Stdout),
	}
	logger.Info("📈 Welcome to Trade 📈")

	market := service.NewMarket(logger)

	var wg sync.WaitGroup
	wg.Add(2)

	go market.Run()

	httpSrv := httpsrv.NewHTTPServer(logger)

	go httpSrv.Run()

	wg.Wait()
}
