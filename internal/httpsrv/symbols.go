package httpsrv

import (
	"net/http"
	"time"

	"github.com/jesseobrien/trade/internal/types"
	"github.com/labstack/echo/v4"
)

func (h HTTPSrv) GetSymbols(ctx echo.Context) error {

	var getSymbolsResponse []types.Symbol

	if err := h.conn.Request("symbols.get", nil, &getSymbolsResponse, time.Second); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, getSymbolsResponse)
}
