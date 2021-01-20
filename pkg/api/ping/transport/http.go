package transport

import (
	"log"
	"net/http"

	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/ping"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo/v4"
)

// HTTP represents password http transport service
type HTTP struct {
	svc ping.Service
}

func NewHTTP(svc ping.Service, er *echo.Group) {
	h := HTTP{svc}
	//	er.POST()
	//	pr := er.Group("/unauthping")
	er.GET("/unauthping", h.createPing)
	er.GET("/authping", h.createPing, authMw.Authenticate())

}

func (h *HTTP) createPing(c echo.Context) error {
	log.Println("Inside createPing()")
	_ = h.svc.Create(c, authapi.Ping{})
	return c.String(http.StatusOK, "pong")

}
