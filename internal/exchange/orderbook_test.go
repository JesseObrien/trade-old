package exchange

import (
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/shopspring/decimal"
	"gopkg.in/go-playground/assert.v1"
)

func TestOrderBookFillsBuysOnInsert(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testSell := orders.New(symbol)
	testSell.Side = orders.SELLSIDE
	testSell.Type = orders.LIMIT
	testSell.Quantity = decimal.NewFromInt(10)
	testSell.Price = decimal.NewFromInt(1)

	ob.Insert(testSell)

	assert.Equal(t, ob.Offers.Len(), 1)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BUYSIDE
	testBuy.Type = orders.LIMIT
	testBuy.Quantity = decimal.NewFromInt(10)
	testBuy.Price = decimal.NewFromInt(1)

	ob.Insert(testBuy)

	assert.Equal(t, ob.Bids.Len(), 0)
	assert.Equal(t, ob.Offers.Len(), 0)
}

func TestOrderBookFillsSellsOnInsert(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BUYSIDE
	testBuy.Type = orders.LIMIT
	testBuy.Quantity = decimal.NewFromInt(10)
	testBuy.Price = decimal.NewFromInt(1)

	ob.Insert(testBuy)
	assert.Equal(t, ob.Bids.Len(), 1)

	testSell := orders.New(symbol)
	testSell.Side = orders.SELLSIDE
	testSell.Type = orders.LIMIT
	testSell.Quantity = decimal.NewFromInt(10)
	testSell.Price = decimal.NewFromInt(1)

	ob.Insert(testSell)

	assert.Equal(t, ob.Bids.Len(), 0)
	assert.Equal(t, ob.Offers.Len(), 0)
}
