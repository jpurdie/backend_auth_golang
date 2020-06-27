package database

import (
	"github.com/go-pg/pg"
	"github.com/jpurdie/authapi"
)

type AuthorizationStore struct {
	db *pg.DB
}

func NewAuthorizationStore(db *pg.DB) *AuthorizationStore {
	return &AuthorizationStore{
		db: db,
	}
}

func (s *AuthorizationStore) CheckAuthorization(u *authapi.User) error {
	return nil
}
