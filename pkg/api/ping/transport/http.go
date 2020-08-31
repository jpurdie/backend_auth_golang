package transport

import (
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/ping"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/labstack/echo"
	"log"
	"net/http"
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
	authMiddleware := authMw.Authenticate()
	er.Use(authMiddleware)
	er.GET("/authping", h.createPing)

}


func (h *HTTP) createPing(c echo.Context) error {
	log.Println("Inside createPing()")
	_ = h.svc.Create(c, authapi.Ping{})
	return c.String(http.StatusOK, "pong")

}
