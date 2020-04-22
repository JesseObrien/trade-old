package orders

import (
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
	ID                   string          `json:"-"`
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

func (o *Order) IsClosed() bool {
	return o.OpenQuantity().Equals(decimal.Zero)
}

func (o *Order) OpenQuantity() decimal.Decimal {
	if o.openQuantity == nil {
		return o.Quantity.Sub(o.ExecutedQuantity)
	}

	return *o.openQuantity
}

func (o *Order) Execute(price, quantity decimal.Decimal) {
	o.ExecutedQuantity = quantity
	o.LastExecutedPrice = price
	o.LastExecutedQuantity = quantity

}

func (o *Order) Cancel() {
	openQuantity := decimal.Zero
	o.openQuantity = &openQuantity
}
