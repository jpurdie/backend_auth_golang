package database

import (
	"github.com/go-pg/pg"
)

type MiddlewareStore struct {
	db *pg.DB
}

func NewMiddlewareStore(db *pg.DB) *MiddlewareStore {
	return &MiddlewareStore{
		db: db,
	}
}

func (s *MiddlewareStore) Check() error {
	return nil
}
