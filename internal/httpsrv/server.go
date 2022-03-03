package httpsrv

import (
	"github.com/apex/log"
	"github.com/labstack/echo/v4"

	nats "github.com/nats-io/nats.go"
)

// HTTPSrv is a new HTTP Server
type HTTPSrv struct {
	logger log.Logger
	conn   *nats.EncodedConn
}

// NewHTTPServer factories out a new HTTP Server
func NewHTTPServer(logger log.Logger, conn *nats.EncodedConn) *HTTPSrv {
	return &HTTPSrv{
		logger,
		conn,
	}
}

// Run starts the HTTP server and binds all handlers
func (h *HTTPSrv) Run() {
	e := echo.New()
	e.HideBanner = true

	e.GET("/symbols", h.GetSymbols)

	e.POST("/orders", h.NewOrder)
	e.DELETE("/orders", h.CancelOrder)

	e.Logger.Fatal(e.Start(":8088"))
}
