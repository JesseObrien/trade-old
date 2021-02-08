package exchange

import (
	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/orders"
)

func (ex *Exchange) HandleNewOrders() {

	newOrdersChan := make(chan *orders.Order)
	sub, err := ex.natsConn.BindRecvChan("order.created", newOrdersChan)
	if err != nil {
		ex.logger.Errorf("Could not bind nats connection for order.created: %v", err)
	}

	for {
		select {
		case order := <-newOrdersChan:
			go ex.onNewOrder(order)
		case <-ex.quit:
			sub.Unsubscribe()
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

	orderbook, ok := ex.Symbols[order.Symbol]
	if !ok {
		ex.logger.WithFields(log.Fields{
			"Symbol": order.Symbol,
		}).Error("symbol is not registered with the exchange")
		return
	}

	orderbook.Insert(order)

	ex.logger.Info(orderbook.Report())
}
