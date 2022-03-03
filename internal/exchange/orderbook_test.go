package exchange

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/jesseobrien/trade/internal/types"
	"github.com/shopspring/decimal"
)

func TestOrderBook_FillsMarketBuysOnInsert(t *testing.T) {
	symbol := types.Symbol("TEST")

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.MarketOrder
	testSell.Quantity = decimal.NewFromInt(10)
	testSell.Price = decimal.NewFromInt(1)

	ob.Insert(testSell)

	assert.Equal(t, ob.Offers.Len(), 1)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.MarketOrder
	testBuy.Quantity = decimal.NewFromInt(10)
	testBuy.Price = decimal.NewFromInt(1)

	ob.Insert(testBuy)

	assert.True(t, testBuy.IsClosed())
	assert.True(t, testSell.IsClosed())
}

func TestOrderBook_FillsMarketSellsOnInsert(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.MarketOrder
	testBuy.Quantity = decimal.NewFromInt(10)
	testBuy.Price = decimal.NewFromInt(1)

	ob.Insert(testBuy)
	assert.Equal(t, ob.Bids.Len(), 1)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.MarketOrder
	testSell.Quantity = decimal.NewFromInt(10)
	testSell.Price = decimal.NewFromInt(1)

	ob.Insert(testSell)

	assert.True(t, testSell.IsClosed())
	assert.True(t, testBuy.IsClosed())
}

func TestOrderBook_FillsLimitBuyWhenLimitExists(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.LimitOrder
	testBuy.Quantity = decimal.NewFromInt(10)

	ob.Insert(testBuy)
	assert.Equal(t, ob.Bids.Len(), 1)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.MarketOrder
	testSell.Quantity = decimal.NewFromInt(10)

	// Set the buy price to $1.00 and the sell price to $0.99
	// It should execute because the buy price is higher
	testBuy.Price = decimal.NewFromFloat(1.00)
	testSell.Price = decimal.NewFromFloat(0.99)

	ob.Insert(testSell)

	assert.True(t, testSell.IsClosed())
	assert.True(t, testBuy.IsClosed())
}

func TestOrderBook_DoesNotFillLimitBuyWhenLimitExceeded(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.LimitOrder
	testBuy.Quantity = decimal.NewFromInt(10)

	ob.Insert(testBuy)
	assert.Equal(t, ob.Bids.Len(), 1)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.MarketOrder
	testSell.Quantity = decimal.NewFromInt(10)

	// Set the buy price to $1.00 and the sell price to $1.01
	// It should not execute because the buy price is lower
	testBuy.Price = decimal.NewFromFloat(1.00)
	testSell.Price = decimal.NewFromFloat(1.01)

	ob.Insert(testSell)

	assert.True(t, testSell.IsClosed())
	assert.True(t, testBuy.IsClosed())
}

func TestOrderBook_FIFOBuyingExecutesCorrectSellOrder(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.MarketOrder
	testBuy.Quantity = decimal.NewFromInt(10)

	testExecuteSell := orders.New(symbol)
	testExecuteSell.Side = orders.SellSide
	testExecuteSell.Type = orders.LimitOrder
	testExecuteSell.Quantity = decimal.NewFromInt(10)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.MarketOrder
	testSell.Quantity = decimal.NewFromInt(10)

	// Set the buy price to $1.00 and the sell price to $0.99
	// It should execute because the buy price is higher
	testBuy.Price = decimal.NewFromFloat(1.00)
	testSell.Price = decimal.NewFromFloat(0.99)
	// The execute order and previous price are identical, the execute should win as it's the first in
	testExecuteSell.Price = decimal.NewFromFloat(0.99)

	ob.Insert(testExecuteSell)
	ob.Insert(testSell)

	ob.Insert(testBuy)

	assert.False(t, testSell.IsClosed())
	assert.True(t, testExecuteSell.IsClosed())
	assert.True(t, testBuy.IsClosed())
}

func TestOrderBook_FillsLimitSellWhenLimitExists(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.MarketOrder
	testBuy.Quantity = decimal.NewFromInt(10)

	ob.Insert(testBuy)
	assert.Equal(t, ob.Bids.Len(), 1)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.LimitOrder
	testSell.Quantity = decimal.NewFromInt(10)

	// Set the buy price to $1.00 and the sell price to $0.99
	// It should execute because the buy price is higher
	testBuy.Price = decimal.NewFromFloat(1.00)
	testSell.Price = decimal.NewFromFloat(0.99)

	ob.Insert(testSell)

	assert.True(t, testSell.IsClosed())
	assert.True(t, testBuy.IsClosed())
}

func TestOrderBook_DoesNotFillLimitSellWhenLimitExceeded(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.MarketOrder
	testBuy.Quantity = decimal.NewFromInt(10)

	ob.Insert(testBuy)
	assert.Equal(t, ob.Bids.Len(), 1)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.LimitOrder
	testSell.Quantity = decimal.NewFromInt(10)

	// Set a limit sell order for $0.99, it should not execute unless the buy price is higher
	testBuy.Price = decimal.NewFromFloat(0.10)
	testSell.Price = decimal.NewFromFloat(0.99)

	ob.Insert(testSell)

	assert.Equal(t, ob.Bids.Len(), 1)
	assert.Equal(t, ob.Offers.Len(), 1)
	assert.False(t, testSell.IsClosed())
}

func TestOrderBook_AveragePriceSellIsAccurate(t *testing.T) {
	symbol := "TEST"

	ob := NewOrderBook(log.Logger{
		Handler: memory.New(),
		Level:   log.DebugLevel,
	}, symbol)

	testBuy := orders.New(symbol)
	testBuy.Side = orders.BuySide
	testBuy.Type = orders.MarketOrder
	testBuy.Quantity = decimal.NewFromFloat(10)
	testBuy.Price = decimal.NewFromFloat(1.00)

	testBuyTwo := orders.New(symbol)
	testBuyTwo.Side = orders.BuySide
	testBuyTwo.Type = orders.MarketOrder
	testBuyTwo.Quantity = decimal.NewFromInt(10)
	testBuyTwo.Price = decimal.NewFromFloat(0.90)

	ob.Insert(testBuy)
	ob.Insert(testBuyTwo)

	testSell := orders.New(symbol)
	testSell.Side = orders.SellSide
	testSell.Type = orders.LimitOrder
	testSell.Quantity = decimal.NewFromInt(20)

	// Set a limit sell order for $0.99, it should not execute unless the buy price is higher
	testSell.Price = decimal.NewFromFloat(0.85)

	ob.Insert(testSell)

	assert.True(t, testBuy.IsClosed())
	assert.True(t, testBuyTwo.IsClosed())
	assert.True(t, testSell.IsClosed())
	t.Logf("avg price %v", testSell.AvgPrice)
	assert.Equal(t, 0, testSell.AvgPrice.Cmp(decimal.NewFromFloat(0.95)))

}
