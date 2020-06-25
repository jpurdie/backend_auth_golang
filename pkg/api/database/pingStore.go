package database

import (
	"github.com/go-pg/pg"
)

type PingStore struct {
	db *pg.DB
}

// NewAdmAccountStore returns an AccountStore.
func NewPingStore(db *pg.DB) *PingStore {
	return &PingStore{
		db: db,
	}
}
func (s *PingStore) Ping() error {
	return nil
}
