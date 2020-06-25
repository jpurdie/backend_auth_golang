package api

import (
	"github.com/jpurdie/authapi/pkg/api/app"
	"github.com/jpurdie/authapi/pkg/api/public"
	"github.com/jpurdie/authapi/pkg/utl/config"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/jpurdie/authapi/pkg/utl/server"
	"github.com/labstack/echo/middleware"
	"log"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	log.Println("Inside Start()")

	db, _ := postgres.DBConn()

	e := server.New()
	e.Pre(middleware.RemoveTrailingSlash())
	//	e.Use(secure.Headers())
	//	e.Use(secure.CORS())

	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	publicAPI, err := public.NewAPI(db)

	publicG := e.Group("/public")
	publicAPI.Router(publicG)

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
