package mock

import (
	"github.com/jpurdie/authapi"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(authapi.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u authapi.User) (string, error) {
	return j.GenerateTokenFn(u)
}
