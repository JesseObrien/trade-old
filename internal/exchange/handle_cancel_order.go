package exchange

import (
	"fmt"

	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/types/requests"
)

// HandleCancelOrderRequest processes cancel order requests
func (ex *Exchange) HandleCancelOrderRequest() {
	sub, err := ex.natsConn.Subscribe("order.cancelled", func(subj, replySubject string, request requests.CancelOrder) {
		ex.logger.WithFields(log.Fields{
			"Symbol":  request.Symbol,
			"OrderID": request.OrderID,
		}).Infof("received an order cancellation request")

		var message string
		var cancelled bool

		orderbook, ok := ex.OrderBooks[request.Symbol]

		if !ok {
			message = fmt.Sprintf("no orderbook found for symbol %s", request.Symbol)
		} else {
			cancelledOrder := orderbook.Cancel(request.OrderID)
			cancelled = true

			if cancelledOrder == nil {
				message = "no order found"
				cancelled = false
			}
		}

		ex.logger.Info(orderbook.Report())

		err := ex.natsConn.Publish(replySubject, &requests.CancelOrderResponse{
			Request:   request,
			Cancelled: cancelled,
			Message:   message,
		})
		if err != nil {
			panic(err)
		}
	})

	if err != nil {
		panic(err)
	}

	for {
		<-ex.quit
		err := sub.Unsubscribe()
		if err != nil {
			panic(err)
		}
		return
	}

}
