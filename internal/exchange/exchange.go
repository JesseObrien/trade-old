package exchange

import (
	"os"
	"os/signal"

	"github.com/jesseobrien/trade/internal/orders"

	"github.com/apex/log"
	"github.com/nats-io/nats.go"

	"github.com/jesseobrien/trade/internal/market"
	"github.com/shopspring/decimal"
)

type Exchange struct {
	logger   log.Logger
	Symbols  map[string]*market.Market
	quit     chan os.Signal
	natsConn *nats.EncodedConn
}

func New(logger log.Logger, conn *nats.EncodedConn) *Exchange {
	return &Exchange{
		logger:   logger,
		Symbols:  map[string]*market.Market{},
		natsConn: conn,
	}
}

func (ex *Exchange) Run() {
	ex.quit = make(chan os.Signal, 1)

	signal.Notify(ex.quit, os.Interrupt)

	defer signal.Stop(ex.quit)

	go ex.HandleNewOrders()
	go ex.HandleCancelOrderRequest()

	m := market.New(ex.logger, "JOBR")
	price, _ := decimal.NewFromString("2.00")
	ex.IPO(m, price, 10000)

	<-ex.quit
	ex.logger.Info("⏳ Shutting down...")
}

// Stop will close the exchange channel
func (ex *Exchange) Stop() {
	close(ex.quit)
}

func (ex *Exchange) IPO(m *market.Market, price decimal.Decimal, sharesIssued int64) {
	quantityShares := decimal.NewFromInt(sharesIssued)
	marketCap := price.Mul(quantityShares)
	ex.logger.Infof("⚡ New Company IPO: %s issuing %d shares @ $%s/share. Market Cap: $%s", m.Symbol, sharesIssued, price.StringFixed(2), marketCap.StringFixedBank(2))

	ex.Symbols[m.Symbol] = m

	o := orders.New(m.Symbol)
	o.Quantity = quantityShares
	o.Price = price
	o.Side = orders.SELLSIDE
	o.Type = orders.MARKET

	m.Insert(o)

	ex.logger.Info(m.Report())
}
