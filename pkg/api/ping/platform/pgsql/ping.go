package pgsql

import (
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi"
	"log"
)

// Custom errors
//var (
//	ErrCompAlreadyExists  = errors.New("Organization name already exists")
//	ErrEmailAlreadyExists = errors.New("Email already exists")
//)

type Ping struct{}


// Create creates a new user on database
func (p Ping) Create(db sqlx.DB, profile authapi.Ping) error {
	log.Println("Ping creation was successful")
	return nil
}
