package requests

// CancelOrder allows us to cancel an order
type CancelOrder struct {
	Symbol  string
	OrderID string `json:"order_id"`
}

// CancelOrderResponse responds to a CancelOrder Request
type CancelOrderResponse struct {
	Request   CancelOrder
	Cancelled bool
	Message   string
}
