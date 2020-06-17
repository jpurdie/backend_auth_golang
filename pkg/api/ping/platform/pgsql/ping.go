package pgsql

import "github.com/jpurdie/authapi"

type Pong struct{}

func (p Pong) Create(id int) (authapi.Pong, error) {
	var pong authapi.Pong
	pong.Resp = "pong"
	return pong, nil
}
