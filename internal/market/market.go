package market

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/orders"
)

type Market struct {
	logger     log.Logger
	Symbol     string
	Bids       *orders.OrderList
	Offers     *orders.OrderList
	Executions []Execution

	mu sync.Mutex
}

func New(logger log.Logger, symbol string) *Market {
	return &Market{
		logger: logger,
		Symbol: symbol,
		Bids:   &orders.OrderList{},
		Offers: &orders.OrderList{},
	}
}

const displayMarketOrders = `
# Market Report for **{{.Symbol}}**

| Open Bids |---|---|
| ID | Price | Quantity |
{{- .Bids.Display -}}

---

| Open Offers |---|---|
| ID | Price | Quantity |
{{- .Offers.Display -}}
`

// Report writes out the current market report for the symbol
func (m *Market) Report() string {
	var buf bytes.Buffer
	t := template.Must(template.New("market").Parse(displayMarketOrders))

	err := t.Execute(&buf, struct {
		Symbol string
		Bids   *orders.OrderList
		Offers *orders.OrderList
	}{
		m.Symbol,
		m.Bids,
		m.Offers,
	})

	if err != nil {
		panic(err)
	}

	return buf.String()
}

// Insert puts an order into the orderlist
func (m *Market) Insert(order *orders.Order) {
	if order.Side == orders.BUYSIDE {
		m.Bids.Insert(order)
	} else {
		m.Offers.Insert(order)
	}
}

// Cancel will cancel an order
func (m *Market) Cancel(oid string) (order *orders.Order) {
	order = m.Bids.Remove(oid)

	if order == nil {
		order = m.Offers.Remove(oid)
	}

	if order != nil {
		order.Cancel()
	}

	return
}

func (m *Market) Match() (matches []orders.Order) {
	// @TODO this method works for matching market orders, need to figure out a way to match others

	for m.Bids.Len() > 0 && m.Offers.Len() > 0 {
		m.mu.Lock()
		bestBid := m.Bids.GetBest()
		bestOffer := m.Offers.GetBest()

		price := bestOffer.Price
		quantity := bestBid.OpenQuantity()

		// @TODO figure out limit orders
		// if bestBid.Type == orders.LIMIT && bestBid.Price.Cmp(bestOffer.Price) {
		// 	x
		// }

		if offerQuantity := bestOffer.OpenQuantity(); offerQuantity.Cmp(quantity) == -1 {
			quantity = offerQuantity
		}

		// @TODO probably need mutex locks while we match orders and then possibly remove them
		bestBid.Execute(price, quantity)
		bestOffer.Execute(price, quantity)

		matches = append(matches, *bestBid, *bestOffer)

		if bestBid.IsClosed() {
			m.Bids.Remove(bestBid.ID)
		}

		if bestOffer.IsClosed() {
			m.Offers.Remove(bestOffer.ID)
		}
		m.mu.Unlock()
	}

	return
}
