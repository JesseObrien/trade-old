package exchange

import (
	"github.com/google/uuid"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/shopspring/decimal"
)

type ExecutionType string

const (
	FILLED           ExecutionType = "FILLED"
	PARTIALLY_FILLED ExecutionType = "PARTIALLY_FILLED"
)

type Execution struct {
	TargetFirmID  string
	SendingFirmID string

	OrderID            string
	ExecutionID        string
	ExecutionType      ExecutionType
	Symbol             string
	Side               orders.OrderSide
	LeavesQuantity     decimal.Decimal
	CumulativeQuantity decimal.Decimal
	AvgPrice           decimal.Decimal
	Quantity           decimal.Decimal
	OrderQuantity      decimal.Decimal
	LastShares         decimal.Decimal
	LastPrice          decimal.Decimal
}

func NewExecution(order orders.Order, status ExecutionType) *Execution {
	return &Execution{
		OrderID:            order.ID,
		ExecutionID:        uuid.New().String(),
		ExecutionType:      status,
		Symbol:             order.Symbol,
		Side:               order.Side,
		LeavesQuantity:     order.OpenQuantity(),
		CumulativeQuantity: order.ExecutedQuantity,
		Quantity:           order.Quantity,
		LastShares:         order.LastExecutedQuantity,
		LastPrice:          order.LastExecutedPrice,
	}
}
