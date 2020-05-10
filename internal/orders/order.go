package orders

import (
	"bytes"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type OrderType string

const (
	MARKET OrderType = "MARKET_ORDER"
	LIMIT  OrderType = "LIMIT_ORDER"
)

type OrderSide string

const (
	BUYSIDE  OrderSide = "BUY"
	SELLSIDE OrderSide = "SELL"
)

type Order struct {
	ID                   string          `json:"id"`
	TargetFirmID         string          `json:"to,omitempty"`
	SendingFirmID        string          `json:"from,omitempty"`
	Symbol               string          `json:"symbol"`
	Type                 OrderType       `json:"order_type"`
	Price                decimal.Decimal `json:"price,omitempty"`
	Side                 OrderSide       `json:"order_side"`
	Quantity             decimal.Decimal `json:"quantity"`
	LastExecutedPrice    decimal.Decimal
	ExecutedQuantity     decimal.Decimal
	LastExecutedQuantity decimal.Decimal
	openQuantity         *decimal.Decimal
	insertedAt           time.Time
	AvgPrice             decimal.Decimal
}

func New(symbol string) *Order {
	return &Order{
		ID:     uuid.New().String(),
		Symbol: symbol,
	}
}

func NewMarketOrder(symbol string, quantity int64) (order *Order) {
	order = New(symbol)
	order.Type = MARKET
	order.Quantity = decimal.NewFromInt(quantity)

	return
}

var displayOrderInTable = `
| {{.ID}} | {{.Price}} | {{.Quantity -}} |
`

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

func (o *Order) IsClosed() bool {
	return o.OpenQuantity().Equals(decimal.Zero)
}

func (o *Order) OpenQuantity() decimal.Decimal {
	if o.openQuantity == nil {
		return o.Quantity.Sub(o.ExecutedQuantity)
	}

	return *o.openQuantity
}

// Execute executes a price and quantity update on an order
func (o *Order) Execute(price, quantity decimal.Decimal) {
	o.ExecutedQuantity = o.ExecutedQuantity.Add(quantity)
	o.LastExecutedPrice = price
	o.LastExecutedQuantity = quantity

}

// Cancel cancels an order
func (o *Order) Cancel() {
	openQuantity := decimal.Zero
	o.openQuantity = &openQuantity
}
