package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-pg/pg/v9"
	// DB adapter
	_ "github.com/lib/pq"
)

type dbLogger struct{}

// BeforeQuery hooks before pg queries
func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

// AfterQuery hooks after pg queries
func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	query, err := q.FormattedQuery()
	fmt.Println(query)
	return err
}

// New creates new database connection to a postgres database
func Init() (*pg.DB, error) {

	opt, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db := pg.Connect(opt)
	db = db.WithTimeout(time.Second * time.Duration(5))

	if true {
		db.AddQueryHook(dbLogger{})
	}

	return db, err
}
