package httpsrv

import (
	"net/http"

	"github.com/apex/log"
	"github.com/labstack/echo/v4"
)

type HttpSrv struct {
	logger log.Logger
}

func NewHTTPServer(logger log.Logger) *HttpSrv {
	return &HttpSrv{
		logger,
	}
}

func (h *HttpSrv) Run() {
	e := echo.New()

	e.POST("/orders", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusAccepted, nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
