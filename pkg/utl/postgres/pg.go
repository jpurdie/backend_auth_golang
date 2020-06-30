package postgres

import (
	"log"
	"os"
	"time"

	"github.com/go-pg/pg"
	// DB adapter
	_ "github.com/lib/pq"
)

type logSQL struct{}

func (l *logSQL) BeforeQuery(e *pg.QueryEvent) {}

func (l *logSQL) AfterQuery(e *pg.QueryEvent) {
	query, err := e.FormattedQuery()
	if err != nil {
		panic(err)
	}
	log.Println(query)
}

// New creates new database connection to a postgres database
func DBConn() (*pg.DB, error) {

	opt, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println(err)
		panic(err)
	}
	opt.PoolSize = 3
	opt.PoolTimeout = time.Second * 5
	opt.IdleCheckFrequency = time.Second * 10
	db := pg.Connect(opt)
	db = db.WithTimeout(time.Second * time.Duration(5))
	if err := checkConn(db); err != nil {
		log.Println(err)
		return nil, err
	}
	if true {
		db.AddQueryHook(&logSQL{})
	}

	return db, err
}
func checkConn(db *pg.DB) error {
	var n int
	_, err := db.QueryOne(pg.Scan(&n), "SELECT 1")
	return err
}
