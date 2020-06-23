// Copyright 2017 Emir Ribic. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// GORSK - Go(lang) restful starter kit
//
// API Docs for GORSK v1
//
// 	 Terms Of Service:  N/A
//     Schemes: http
//     Version: 2.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Emir Ribic <ribice@gmail.com> https://ribice.ba
//     Host: localhost:8080
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer: []
//
//     SecurityDefinitions:
//     bearer:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package api

import (
	"github.com/jpurdie/authapi/pkg/api/app"
	"github.com/jpurdie/authapi/pkg/api/public"
	"github.com/jpurdie/authapi/pkg/utl/config"
	"github.com/jpurdie/authapi/pkg/utl/middleware/secure"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/jpurdie/authapi/pkg/utl/server"
	"github.com/labstack/echo/middleware"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {

	db, _ := postgres.DBConn()

	e := server.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(secure.Headers())
	e.Use(secure.CORS())

	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	publicAPI, err := public.NewAPI(db)

	appAPI, err := app.NewAPI(db)
	if err != nil {
		panic(err)
	}

	public := e.Group("/public")
	publicAPI.Router(public)
	api := e.Group("/api")
	v1 := api.Group("/v1")
	appAPI.Router(v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
