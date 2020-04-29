package httpsrv

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jesseobrien/trade/internal/orders"
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

	return ctx.JSON(http.StatusAccepted, nil)

}
