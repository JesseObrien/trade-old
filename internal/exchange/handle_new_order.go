package exchange

import (
	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/orders"
)

func (exch *Exchange) HandleNewOrders() {

	newOrdersChan := make(chan *orders.Order)
	exch.natsConn.BindRecvChan("order.created", newOrdersChan)

	for {
		select {
		case order := <-newOrdersChan:
			go exch.onNewOrder(order)
		case <-exch.quit:
			return
		}
	}
}

func (ex *Exchange) onNewOrder(order *orders.Order) {
	ex.logger.WithFields(log.Fields{
		"Symbol":    order.Symbol,
		"ID":        order.ID,
		"OrderType": order.Type,
		"Side":      order.Side,
		"Quantity":  order.Quantity,
		"Price":     order.Price.StringFixed(2),
		"Value":     order.Price.Mul(order.Quantity).StringFixed(2),
	}).Info("💸 A new order was received!")

	market, ok := ex.Symbols[order.Symbol]
	if !ok {
		ex.logger.WithFields(log.Fields{
			"Symbol": order.Symbol,
		}).Error("symbol is not registered with the exchange")
		return
	}
	market.Insert(order)

	matches := market.Match()

	for _, m := range matches {
		ex.logger.WithFields(log.Fields{
			"Symbol":        m.Symbol,
			"OrderType":     m.Type,
			"Side":          m.Side,
			"QtyFilled":     m.ExecutedQuantity,
			"ExecutedPrice": m.LastExecutedPrice.StringFixed(2),
			"Value":         m.LastExecutedPrice.Mul(order.ExecutedQuantity).StringFixed(2),
		}).Info("Order Executed")
	}

	ex.logger.Info(market.Report())
}

func (ex *Exchange) Match(o *orders.Order) {

}
