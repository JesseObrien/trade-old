package httpsrv

import (
	"net/http"

	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/exchange"
	"github.com/jesseobrien/trade/internal/orders"
	"github.com/labstack/echo/v4"
)

type HttpSrv struct {
	logger   log.Logger
	exchange *exchange.Exchange
}

func NewHTTPServer(logger log.Logger, exchange *exchange.Exchange) *HttpSrv {
	return &HttpSrv{
		logger,
		exchange,
	}
}

func (h *HttpSrv) Run() {
	e := echo.New()

	e.POST("/orders", func(ctx echo.Context) error {
		order := &orders.Order{}

		if err := ctx.Bind(order); err != nil {
			return ctx.JSON(http.StatusBadRequest, err)
		}

		h.exchange.SubmitOrder(order)

		return ctx.JSON(http.StatusAccepted, nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
