package httpsrv

import (
	"net/http"

	"github.com/apex/log"
	"github.com/jesseobrien/trade/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type HttpSrv struct {
	logger      log.Logger
	traderStore *service.TraderStore
}

func NewHTTPServer(logger log.Logger, traderStore *service.TraderStore) *HttpSrv {
	return &HttpSrv{
		logger,
		traderStore,
	}
}

func (h *HttpSrv) Run() {
	e := echo.New()

	e.POST("/traders", func(c echo.Context) error {

		trader := &service.Trader{}

		if err := c.Bind(trader); err != nil {
			err = errors.Wrapf(err, "Could not bind request body")
			h.logger.Error(err.Error())
			return c.JSON(http.StatusBadRequest, err)
		}

		if err := h.traderStore.Write(trader); err != nil {
			h.logger.Errorf("error writing to database", err)
			return c.JSON(http.StatusInternalServerError, "database issue")
		}

		return c.JSON(http.StatusCreated, trader)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
