package database

import (
	"github.com/go-pg/pg"
)

type OrganizationStore struct {
	db *pg.DB
}

func NewOrganizationStore(db *pg.DB) *OrganizationStore {
	return &OrganizationStore{
		db: db,
	}
}
