package market

import (
	"fmt"

	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/orders"
)

type Market struct {
	logger     log.Logger
	Symbol     string
	Bids       orders.OrderList
	Offers     orders.OrderList
	Executions []Execution
}

func New(logger log.Logger, symbol string) *Market {
	return &Market{
		logger: logger,
		Symbol: symbol,
		Bids:   orders.OrderList{},
		Offers: orders.OrderList{},
	}
}

func (m *Market) Report() string {

	return fmt.Sprintf(`
	Market Report for %s

	Open Bids:
	%s
	Open Offers:
	%s
`, m.Symbol, m.Bids.Display(), m.Offers.Display())
}

func (m *Market) Insert(order *orders.Order) {
	if order.Side == orders.BUYSIDE {
		m.Bids.Insert(order)
	} else {
		m.Offers.Insert(order)
	}
}

func (m *Market) Cancel(oid string, side orders.OrderSide) (order *orders.Order) {
	if order.Side == orders.BUYSIDE {
		order = m.Bids.Remove(oid)
	} else {
		order = m.Offers.Remove(oid)
	}

	if order != nil {
		order.Cancel()
	}

	return
}

func (m *Market) Match() (matches []*orders.Order) {
	// @TODO this method works for matching market orders, need to figure out a way to match others
	for m.Bids.Len() > 0 && m.Offers.Len() > 0 {
		bestBid := m.Bids.GetBest()
		bestOffer := m.Offers.GetBest()

		price := bestOffer.Price
		quantity := bestBid.OpenQuantity()

		if offerQuant := bestOffer.OpenQuantity(); offerQuant.Cmp(quantity) == -1 {
			quantity = offerQuant
		}

		// @TODO probably need mutex locks while we match orders and then possibly remove them
		bestBid.Execute(price, quantity)
		bestOffer.Execute(price, quantity)

		matches = append(matches, bestBid, bestOffer)

		if bestBid.IsClosed() {
			m.Bids.Remove(bestBid.ID)
		}

		if bestOffer.IsClosed() {
			m.Offers.Remove(bestOffer.ID)
		}
	}

	return
}
