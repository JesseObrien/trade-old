package orders

import (
	"bytes"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/jesseobrien/trade/internal/types"
	"github.com/shopspring/decimal"
)

type OrderType string

const (
	MarketOrder OrderType = "MARKET_ORDER"
	LimitOrder  OrderType = "LIMIT_ORDER"
)

type OrderSide string

const (
	BuySide  OrderSide = "BUY"
	SellSide OrderSide = "SELL"
)

type Order struct {
	ID                   string          `json:"id"`
	TargetFirmID         string          `json:"to,omitempty"`
	SendingFirmID        string          `json:"from,omitempty"`
	Symbol               types.Symbol    `json:"symbol"`
	Type                 OrderType       `json:"order_type"`
	Price                decimal.Decimal `json:"price,omitempty"`
	Side                 OrderSide       `json:"order_side"`
	Quantity             decimal.Decimal `json:"quantity"`
	LastExecutedPrice    decimal.Decimal
	ExecutedQuantity     decimal.Decimal
	LastExecutedQuantity decimal.Decimal
	openQuantity         *decimal.Decimal
	insertedAt           time.Time
	AvgPrice             decimal.Decimal `json:"average"`
}

// New initialize a new order with an ID set
func New(symbol types.Symbol) *Order {
	return &Order{
		ID:     uuid.New().String(),
		Symbol: symbol,
	}
}

// NewMarketOrder set up a new market order for a certain quantity
func NewMarketOrder(symbol types.Symbol, quantity int64) (order *Order) {
	order = New(symbol)
	order.Type = MarketOrder
	order.Quantity = decimal.NewFromInt(quantity)

	return
}

var displayOrderInTable = `
| {{.ID}} | {{.Price}} | {{.Quantity -}} |
`

// Display shows a report of the order
func (o *Order) Display() string {
	var buf bytes.Buffer
	t := template.Must(template.New("order").Parse(displayOrderInTable))

	err := t.Execute(&buf, struct {
		ID       string
		Price    string
		Quantity string
	}{
		o.ID,
		o.Price.StringFixed(2),
		o.OpenQuantity().String(),
	})

	if err != nil {
		panic(err)
	}

	return buf.String()
}

// IsClosed If the order's open quantity is empty, it's a closed order
func (o *Order) IsClosed() bool {
	return o.OpenQuantity().Equals(decimal.Zero)
}

// OpenQuantity returns the remaining quantityt that hasn't been filled for this order
func (o *Order) OpenQuantity() decimal.Decimal {
	if o.openQuantity == nil {
		return o.Quantity.Sub(o.ExecutedQuantity)
	}

	return *o.openQuantity
}

// Execute executes a price and quantity update on an order
func (o *Order) Execute(price, quantity decimal.Decimal) {
	o.ExecutedQuantity = o.ExecutedQuantity.Add(quantity)

	o.AvgPrice = price

	o.LastExecutedPrice = price
	o.LastExecutedQuantity = quantity
}

// Cancel cancels an order
func (o *Order) Cancel() {
	openQuantity := decimal.Zero
	o.openQuantity = &openQuantity
}
