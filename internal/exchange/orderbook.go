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
	Symbol     string            `json:"symbol"`
	Bids       *orders.OrderList `json:"bids"`
	Offers     *orders.OrderList `json:"offers"`
	Executions []*Execution      `json:"executions"`

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
	if order.Side == orders.BuySide {
		if !ob.FillBuy(order) {
			ob.Bids.Insert(order)
		}
	} else {
		if !ob.FillSell(order) {
			ob.Offers.Insert(order)
		}
	}
}

// Cancel will cancel an open order
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
// FillBuy will try to fill a buy order
func (ob *OrderBook) FillBuy(buyOrder *orders.Order) bool {
	ob.logger.WithFields(log.Fields{
		"ID":       buyOrder.ID,
		"Price":    buyOrder.Price.StringFixed(2),
		"Quantity": buyOrder.Quantity.String(),
	}).Info("received buy order, attempting to fill")

	for {
		if ob.Offers.Len() == 0 {
			return false
		}

		sellOrder := ob.Offers.GetMin()

		if buyOrder.Type == orders.LimitOrder && buyOrder.Price.Cmp(sellOrder.Price) == -1 {
			ob.logger.WithFields(log.Fields{
				"Sell Price": sellOrder.Price.StringFixed(2),
				"Buy Price":  buyOrder.Price.StringFixed(2),
			}).Info("will not fill order, buy limit exceeded")
			return false
		}

		if buyOrder.OpenQuantity().Cmp(sellOrder.OpenQuantity()) == 1 {
			// @TODO need to ensure the sell order is also not a LIMIT and we're exceeding the limit
			var price decimal.Decimal

			priceDiff := buyOrder.Price.Sub(sellOrder.Price)
			if priceDiff.IsZero() {
				price = buyOrder.Price
			} else {
				price = sellOrder.Price.Add(priceDiff.Div(decimal.NewFromInt(2)))
			}
			quantity := sellOrder.OpenQuantity()

			buyOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*buyOrder, PARTIALLY_FILLED))

			sellOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*sellOrder, FILLED))

			ob.Offers.Remove(sellOrder.ID)

			ob.logger.WithFields(log.Fields{
				"ID":       sellOrder.ID,
				"Price":    price.StringFixed(2),
				"Quantity": sellOrder.Quantity.String(),
			}).Info("matched partial sell")

			continue
		}

		if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 1 {
			var price decimal.Decimal

			priceDiff := buyOrder.Price.Sub(sellOrder.Price)
			if priceDiff.IsZero() {
				price = buyOrder.Price
			} else {
				price = sellOrder.Price.Add(priceDiff.Div(decimal.NewFromInt(2)))
			}

			quantity := buyOrder.OpenQuantity()

			buyOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*buyOrder, FILLED))

			sellOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*sellOrder, PARTIALLY_FILLED))

			ob.logger.WithFields(log.Fields{
				"ID":       sellOrder.ID,
				"Price":    price.StringFixed(2),
				"Quantity": sellOrder.Quantity.String(),
			}).Info("matched full sell")

			return true
		}

		if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 0 {
			var price decimal.Decimal

			priceDiff := buyOrder.Price.Sub(sellOrder.Price)
			if priceDiff.IsZero() {
				price = buyOrder.Price
			} else {
				price = sellOrder.Price.Add(priceDiff.Div(decimal.NewFromInt(2)))
			}
			quantity := buyOrder.OpenQuantity()

			buyOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*buyOrder, FILLED))
			sellOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*sellOrder, FILLED))

			ob.Offers.Remove(sellOrder.ID)

			ob.logger.WithFields(log.Fields{
				"ID":       sellOrder.ID,
				"Price":    price.StringFixed(2),
				"Quantity": sellOrder.Quantity.String(),
			}).Info("matched exact sell")

			return true
		}

		ob.logger.Info("could not fill buy")

		return false
	}
}

// FillSell attempts to fill a sell order
func (ob *OrderBook) FillSell(sellOrder *orders.Order) bool {
	ob.logger.WithFields(log.Fields{
		"ID":       sellOrder.ID,
		"Price":    sellOrder.Price.StringFixed(2),
		"Quantity": sellOrder.Quantity.String(),
	}).Info("received sell order, attempting to fill")

	for {
		if ob.Bids.Len() == 0 {
			return false
		}

		buyOrder := ob.Bids.GetMax()

		if sellOrder.Type == orders.LimitOrder && buyOrder.Price.Cmp(sellOrder.Price) == -1 {
			ob.logger.WithFields(log.Fields{
				"Sell Price": sellOrder.Price.StringFixed(2),
				"Buy Price":  buyOrder.Price.StringFixed(2),
			}).Info("will not fill order, sell limit exceeded")
			return false
		}

		if buyOrder.OpenQuantity().Cmp(sellOrder.OpenQuantity()) == 1 {
			var price decimal.Decimal

			priceDiff := buyOrder.Price.Sub(sellOrder.Price)
			if priceDiff.IsZero() {
				price = buyOrder.Price
			} else {
				price = sellOrder.Price.Add(priceDiff.Div(decimal.NewFromInt(2)))
			}
			quantity := sellOrder.Quantity

			buyOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*buyOrder, PARTIALLY_FILLED))

			sellOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*sellOrder, FILLED))

			ob.logger.WithFields(log.Fields{
				"ID":       sellOrder.ID,
				"Price":    price.StringFixed(2),
				"Quantity": sellOrder.Quantity.String(),
			}).Info("matched full buy, filled sell order")

			return true
		}

		// sell amount is larger
		if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 1 {
			var price decimal.Decimal

			priceDiff := buyOrder.Price.Sub(sellOrder.Price)
			if priceDiff.IsZero() {
				price = buyOrder.Price
			} else {
				price = sellOrder.Price.Add(priceDiff.Div(decimal.NewFromInt(2)))
			}

			quantity := buyOrder.OpenQuantity()

			buyOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*buyOrder, FILLED))
			sellOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*sellOrder, PARTIALLY_FILLED))

			ob.Bids.Remove(buyOrder.ID)

			ob.logger.WithFields(log.Fields{
				"ID":       sellOrder.ID,
				"Price":    price.StringFixed(2),
				"Quantity": sellOrder.Quantity.String(),
			}).Info("matched partial buy, partially filled sell order")

			continue
		}

		// sell and buy are the same quantity
		if sellOrder.OpenQuantity().Cmp(buyOrder.OpenQuantity()) == 0 {
			var price decimal.Decimal

			priceDiff := buyOrder.Price.Sub(sellOrder.Price)
			if priceDiff.IsZero() {
				price = buyOrder.Price
			} else {
				price = sellOrder.Price.Add(priceDiff.Div(decimal.NewFromInt(2)))
			}
			quantity := sellOrder.OpenQuantity()

			buyOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*buyOrder, FILLED))
			sellOrder.Execute(price, quantity)
			ob.Executions = append(ob.Executions, NewExecution(*sellOrder, FILLED))

			ob.Bids.Remove(buyOrder.ID)

			ob.logger.WithFields(log.Fields{
				"ID":       sellOrder.ID,
				"Price":    price.StringFixed(2),
				"Quantity": sellOrder.Quantity.String(),
			}).Info("matched exact buy")

			return true
		}

		ob.logger.Info("could not fill sell")

		return false
	}
}
