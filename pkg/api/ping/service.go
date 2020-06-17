package ping

import (
	"github.com/labstack/echo"
)

type Service interface {
	Create(echo.Context, int) (string, error)
}

func New() Ping {
	return Ping{}
}

// Initialize initalizes Ping application service with defaults
func Initialize() Ping {
	return New()
}

type Ping struct {
}
