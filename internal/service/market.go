package service

import (
	"os"
	"os/signal"
	"reflect"

	"github.com/apex/log"
)

type Market struct {
	logger  log.Logger
	Symbols map[string]*Symbol
	Orders  chan Order
	Trades  chan *Trade
	quit    chan os.Signal
}

type Trade struct {
	ID     string
	Price  int64
	Buyer  string
	Seller string
	Shares int64
}

func (t Trade) Value() int64 {
	return t.Price * t.Shares
}

func NewMarket(logger log.Logger) *Market {
	return &Market{
		logger:  logger,
		Symbols: map[string]*Symbol{},
		Orders:  make(chan Order),
		Trades:  make(chan *Trade),
	}
}

func (m *Market) Run() {
	m.quit = make(chan os.Signal, 1)
	signal.Notify(m.quit, os.Interrupt)

	defer signal.Stop(m.quit)

	s := NewSymbol("JOBR", 200, 1000)
	m.IPO(s)

	trader := NewTrader("Jesse O'Brien", 4000)
	m.logger.Infof("🤑 A new trader appears: %s", trader.Name)

	go m.Process()

	m.Orders <- NewMarketOrder(s.Name, trader, 100)

	<-m.quit
	m.logger.Info("⏳ Shutting down...")
}

func (m *Market) Stop() {
	close(m.quit)
}

func (m *Market) Process() {
	for {
		select {
		case trade := <-m.Trades:
			m.logger.WithFields(&log.Fields{
				"Price":  trade.Price,
				"Seller": trade.Seller,
				"Buyer":  trade.Buyer,
				"Shares": trade.Shares,
				"Value":  trade.Value(),
			}).Info("💵 A trade took place!")

		case order := <-m.Orders:
			m.logger.WithFields(log.Fields{
				"OrderType": reflect.TypeOf(order),
				"Symbol":    order.Symbol(),
				"Trader":    order.Trader().Name,
				"Shares":    order.Shares(),
			}).Info("💸 A new order was received!")

			s, ok := m.Symbols[order.Symbol()]
			if !ok {
				// Handle this error and notify the order failed
				return
			}

			s.Mux.Lock()

			for _, ask := range s.Asks {
				if ask.Shares >= order.Shares() {
					t := &Trade{
						Price:  ask.Price,
						Buyer:  order.Trader().Name,
						Seller: ask.Owner,
						Shares: order.Shares(),
					}

					m.Trades <- t
				}

				ask.Shares = ask.Shares - order.Shares()
			}

			defer s.Mux.Unlock()
		case <-m.quit:
			m.Stop()
			return
		}
	}
}

func (m *Market) IPO(symbol *Symbol) {
	m.logger.Infof("⚡ New Company IPO: %s issuing %d shares @ $%.2f ...", symbol.Name, symbol.IssuedShares, symbol.IPOPriceCurrency())

	m.Symbols[symbol.Name] = symbol
}

func (m *Market) SubmitOrder(order Order) {

}
