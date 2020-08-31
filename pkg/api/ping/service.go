package ping

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/jpurdie/authapi"
	"github.com/labstack/echo"
	"github.com/jpurdie/authapi/pkg/api/ping/platform/pgsql"
)

type Service interface {
	Create(echo.Context, authapi.Ping) error
}

// New creates new password application service
func New(db *pg.DB, pdb PingDB) Ping {
	return Ping{
		db:   db,
		pdb:  pdb,
	}
}

func Initialize(db *pg.DB) Ping {
	return New(db, pgsql.Ping{})
}

type Ping struct {
	db   *pg.DB
	pdb  PingDB
}

type PingDB interface {
	Create(orm.DB, authapi.Ping) error
}

