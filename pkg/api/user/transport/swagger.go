package transport

import (
	"github.com/jpurdie/authapi"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*authapi.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []authapi.User `json:"users"`
		Page  int          `json:"page"`
	}
}
