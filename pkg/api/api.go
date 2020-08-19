package api

import (
	"github.com/jpurdie/authapi/pkg/api/app"
	"github.com/jpurdie/authapi/pkg/utl/config"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/jpurdie/authapi/pkg/utl/server"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {

	db, err := postgres.DBConn()
	if err != nil {
		panic(err)
	}

	e := server.New()

	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	appAPI, err := app.NewAPI(db)
	if err != nil {
		panic(err)
	}

	v1 := e.Group("/api/v1")
	appAPI.Router(v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
