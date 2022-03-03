package requests

import "github.com/jesseobrien/trade/internal/types"

// CancelOrder allows us to cancel an order
type CancelOrder struct {
	Symbol  types.Symbol
	OrderID string `json:"order_id"`
}

// CancelOrderResponse responds to a CancelOrder Request
type CancelOrderResponse struct {
	Request   CancelOrder
	Cancelled bool
	Message   string
}
