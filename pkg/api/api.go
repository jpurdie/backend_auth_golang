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
	company "github.com/jpurdie/authapi/pkg/api/company"
	companyLog "github.com/jpurdie/authapi/pkg/api/company/logging"
	companyTransp "github.com/jpurdie/authapi/pkg/api/company/transport"
	"github.com/jpurdie/authapi/pkg/api/ping"
	pingl "github.com/jpurdie/authapi/pkg/api/ping/logging"
	pingt "github.com/jpurdie/authapi/pkg/api/ping/transport"
	user "github.com/jpurdie/authapi/pkg/api/user"
	userLog "github.com/jpurdie/authapi/pkg/api/user/logging"
	userTransp "github.com/jpurdie/authapi/pkg/api/user/transport"
	"github.com/jpurdie/authapi/pkg/utl/config"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/jpurdie/authapi/pkg/utl/server"
	"github.com/jpurdie/authapi/pkg/utl/zlog"
	"github.com/labstack/echo"
	"net/http"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.Init()
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)
	public := e.Group("public")

	pingt.NewHTTP(pingl.New(ping.Initialize(), log), public)

	v1 := e.Group("api/v1")

	companyTransp.NewHTTP(companyLog.New(company.Initialize(db), log), v1)

	//setting up middleware for authentication
	authMiddleware := authMw.Middleware()
	v1.Use(authMiddleware)

	//test ping pong - behind authentication
	v1.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	userTransp.NewHTTP(userLog.New(user.Initialize(db), log), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
