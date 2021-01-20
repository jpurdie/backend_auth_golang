package ping

import (
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"github.com/jpurdie/authapi/pkg/api/ping/platform/pgsql"
	"github.com/labstack/echo/v4"
)

type Service interface {
	Create(echo.Context, authapi.Ping) error
}

// New creates new password application service
func New(db *sqlx.DB, pdb PingDB) Ping {
	return Ping{
		db:  db,
		pdb: pdb,
	}
}

func Initialize(db *sqlx.DB) Ping {
	return New(db, pgsql.Ping{})
}

type Ping struct {
	db  *sqlx.DB
	pdb PingDB
}

type PingDB interface {
	Create(sqlx.DB, authapi.Ping) error
}
