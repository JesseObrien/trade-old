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

	e.POST("/traders", func(c echo.Context) error {
		return c.JSON(http.StatusCreated, nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
