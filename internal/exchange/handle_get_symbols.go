package exchange

import "github.com/nats-io/nats.go"

func (ex *Exchange) HandleGetSymbolsRequest() {

	symbolsRequests := make(chan *nats.Msg, 64)
	sub, err := ex.natsConn.Conn.ChanSubscribe("symbols.get", symbolsRequests)

	if err != nil {
		ex.logger.Errorf("Error: %v", err)
		return
	}

	for {
		select {
		case msg := <-symbolsRequests:
			go ex.onGetSymbols(msg)
		case <-ex.quit:
			sub.Unsubscribe()
			return
		}
	}

}

func (ex *Exchange) onGetSymbols(msg *nats.Msg) {

	ex.logger.Info("Symbols requested")

	if err := ex.natsConn.Publish(msg.Reply, ex.OrderBooks); err != nil {
		ex.logger.Errorf("error replying with symbols: %v", err)
		return
	}

}
