package exchange

import (
	"os"
	"os/signal"

	"github.com/apex/log"

	"github.com/jesseobrien/trade/internal/market"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/shopspring/decimal"
)

type Exchange struct {
	logger  log.Logger
	Symbols map[string]*market.Market
	Orders  chan *orders.Order
	quit    chan os.Signal
}

func New(logger log.Logger) *Exchange {
	return &Exchange{
		logger:  logger,
		Symbols: map[string]*market.Market{},
		Orders:  make(chan *orders.Order),
	}
}

func (ex *Exchange) Run() {
	ex.quit = make(chan os.Signal, 1)

	signal.Notify(ex.quit, os.Interrupt)

	defer signal.Stop(ex.quit)

	go ex.HandleOrders()

	m := market.New(ex.logger, "JOBR")
	price, _ := decimal.NewFromString("2.00")
	ex.IPO(m, price, 10000)

	o := orders.NewMarketOrder(m.Symbol, 100)
	o.Side = orders.BUYSIDE

	ex.Orders <- o

	<-ex.quit
	ex.logger.Info("⏳ Shutting down...")
}

func (ex *Exchange) Stop() {
	close(ex.quit)
}

func (ex *Exchange) IPO(m *market.Market, price decimal.Decimal, sharesIssued int64) {
	quantityShares := decimal.NewFromInt(sharesIssued)
	marketCap := price.Mul(quantityShares)
	ex.logger.Infof("⚡ New Company IPO: %s issuing %d shares @ $%s/share. Market Cap: $%s", m.Symbol, sharesIssued, price.StringFixed(2), marketCap.StringFixedBank(2))

	ex.Symbols[m.Symbol] = m
	// @TODO inject an offer with all of the shares at the price

	o := orders.New(m.Symbol)
	o.Quantity = quantityShares
	o.Price = price
	o.Side = orders.SELLSIDE
	o.Type = orders.MARKET

	ex.Orders <- o
}

func (ex *Exchange) SubmitOrder(orders *orders.Order) {

}
