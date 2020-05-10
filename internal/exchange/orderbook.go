package exchange

import (
	"bytes"
	"sync"
	"text/template"

	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	logger     log.Logger
	Symbol     string
	Bids       *orders.OrderList
	Offers     *orders.OrderList
	Executions []Execution

	MarketPrice decimal.Decimal

	mu sync.Mutex
}

func NewOrderBook(logger log.Logger, symbol string) *OrderBook {
	return &OrderBook{
		logger: logger,
		Symbol: symbol,
		Bids:   &orders.OrderList{},
		Offers: &orders.OrderList{},
	}
}

const displayOrderBookOrders = `
# OrderBook Report for **{{.Symbol}}**

| Open Bids |---|---|
| ID | Price | Quantity |
{{- .Bids.Display -}}

---

| Open Offers |---|---|
| ID | Price | Quantity |
{{- .Offers.Display -}}
`

// Report writes out the current orderbook report for the symbol
func (ob *OrderBook) Report() string {
	var buf bytes.Buffer
	t := template.Must(template.New("orderbook").Parse(displayOrderBookOrders))

	err := t.Execute(&buf, struct {
		Symbol string
		Bids   *orders.OrderList
		Offers *orders.OrderList
	}{
		ob.Symbol,
		ob.Bids,
		ob.Offers,
	})

	if err != nil {
		panic(err)
	}

	return buf.String()
}

// Insert puts an order into the orderlist
func (ob *OrderBook) Insert(order *orders.Order) {
	if order.Side == orders.BUYSIDE {
		if !ob.FillBuy(order) {
			ob.Bids.Insert(order)
		}
	} else {
		if !ob.FillSell(order) {
			ob.Offers.Insert(order)
		}
	}
}

// Cancel will cancel an order
func (ob *OrderBook) Cancel(oid string) (order *orders.Order) {
	order = ob.Bids.Remove(oid)

	if order == nil {
		order = ob.Offers.Remove(oid)
	}

	if order != nil {
		order.Cancel()
	}

	return
}

// Inspiration https://github.com/fmstephe/matching_engine/blob/1f9ff299e0fc65cefd7acd0a637d5764575f4996/matcher/matcher.go

func (ob *OrderBook) FillBuy(buyOrder *orders.Order) bool {
	lwf := ob.logger.WithFields(log.Fields{
		"ID":       buyOrder.ID,
		"Price":    buyOrder.Price.StringFixed(2),
		"Quantity": buyOrder.Quantity.String(),
	})

	lwf.Info("received buy order, attempting to fill")

	for {

		if ob.Offers.Len() == 0 {
			return false
		}

		sellOrder := ob.Offers.GetMin()

		if buyOrder.Price.Cmp(sellOrder.Price) != -1 {
			if buyOrder.OpenQuantity().Cmp(sellOrder.OpenQuantity()) == 1 {
				lwf.Info("matched partial sell")
				price := sellOrder.Price
				quantity := sellOrder.OpenQuantity()

				buyOrder.Execute(price, quantity)
				sellOrder.Execute(price, quantity)

				ob.Offers.Remove(sellOrder.ID)
				continue
			}

			if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 1 {
				lwf.Info("matched full sell")
				price := sellOrder.Price
				quantity := buyOrder.OpenQuantity()

				buyOrder.Execute(price, quantity)
				sellOrder.Execute(price, quantity)

				return true
			}

			if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 0 {
				lwf.Info("matched matched exact sell")
				price := sellOrder.Price
				quantity := buyOrder.OpenQuantity()

				buyOrder.Execute(price, quantity)
				sellOrder.Execute(price, quantity)

				ob.Offers.Remove(sellOrder.ID)

				return true
			}
		}

		lwf.Info("could not fill buy")

		return false
	}
}

func (ob *OrderBook) FillSell(sellOrder *orders.Order) bool {
	lwf := ob.logger.WithFields(log.Fields{
		"ID":       sellOrder.ID,
		"Price":    sellOrder.Price.StringFixed(2),
		"Quantity": sellOrder.Quantity.String(),
	})

	lwf.Info("received sell order, attempting to fill")

	for {

		if ob.Bids.Len() == 0 {
			return false
		}

		buyOrder := ob.Bids.GetMax()

		if buyOrder.Price.Cmp(sellOrder.Price) != -1 {

			if buyOrder.OpenQuantity().Cmp(sellOrder.OpenQuantity()) == 1 {
				lwf.Info("matched full buy")
				price := buyOrder.Price
				quantity := sellOrder.Quantity

				buyOrder.Execute(price, quantity)
				sellOrder.Execute(price, quantity)

				return true
			}

			// sell amount is larger
			if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 1 {
				lwf.Info("matched partial buy")
				price := buyOrder.Price
				quantity := buyOrder.OpenQuantity()

				buyOrder.Execute(price, quantity)
				sellOrder.Execute(price, quantity)

				ob.Bids.Remove(buyOrder.ID)

				continue
			}

			// buy and sell are same
			if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 0 {
				lwf.Info("matched exact buy")
				price := buyOrder.Price
				quantity := sellOrder.OpenQuantity()

				buyOrder.Execute(price, quantity)
				sellOrder.Execute(price, quantity)

				ob.Bids.Remove(buyOrder.ID)
				return true
			}
		}
		lwf.Info("could not fill sell")

		return false
	}
}

// ob.mu.Lock()
// 	for ob.Bids.Len() > 0 && ob.Offers.Len() > 0 {
// 		bestBid := ob.Bids.GetBest()
// 		bestOffer := ob.Offers.GetBest()

// 		price := bestOffer.Price
// 		quantity := bestBid.OpenQuantity()

// 		if offerQuantity := bestOffer.OpenQuantity(); offerQuantity.Cmp(quantity) == -1 {
// 			quantity = offerQuantity
// 		}

// 		// @TODO probably need mutex locks while we match orders and then possibly remove them
// 		bestBid.Execute(price, quantity)
// 		bestOffer.Execute(price, quantity)

// 		matches = append(matches, *bestBid, *bestOffer)

// 		if bestBid.IsClosed() {
// 			ob.Bids.Remove(bestBid.ID)
// 		}

// 		if bestOffer.IsClosed() {
// 			ob.Offers.Remove(bestOffer.ID)
// 		}
// 	}

// 	ob.mu.Unlock()
