package httpsrv

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/jesseobrien/trade/internal/types/requests"
	"github.com/labstack/echo/v4"
)

// NewOrder handles new order creation
func (h HTTPSrv) NewOrder(ctx echo.Context) error {
	order := &orders.Order{}

	if err := ctx.Bind(order); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	order.ID = uuid.New().String()

	newOrdersCh := make(chan *orders.Order)
	h.conn.BindSendChan("order.created", newOrdersCh)

	newOrdersCh <- order

	h.logger.Infof("order.created with id %v", order.ID)

	return ctx.JSON(http.StatusAccepted, struct {
		OrderID string `json:"order_id"`
		Symbol  string
	}{OrderID: order.ID, Symbol: order.Symbol})

}

// CancelOrder puts out cancellation requests
func (h HTTPSrv) CancelOrder(ctx echo.Context) error {
	cancelOrder := &requests.CancelOrder{}

	if err := ctx.Bind(cancelOrder); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	h.logger.Infof("received cancel order request for order id: %s", cancelOrder.OrderID)

	response := &requests.CancelOrderResponse{}
	err := h.conn.Request("order.cancelled", cancelOrder, response, 1*time.Second)

	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, response)
}
