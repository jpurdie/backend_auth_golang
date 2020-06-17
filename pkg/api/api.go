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
	"crypto/sha1"
	company "github.com/jpurdie/authapi/pkg/api/company"
	companyLog "github.com/jpurdie/authapi/pkg/api/company/logging"
	companyTransp "github.com/jpurdie/authapi/pkg/api/company/transport"
	"github.com/jpurdie/authapi/pkg/api/ping"
	pingl "github.com/jpurdie/authapi/pkg/api/ping/logging"
	pingt "github.com/jpurdie/authapi/pkg/api/ping/transport"
	"github.com/jpurdie/authapi/pkg/utl/config"
	"github.com/jpurdie/authapi/pkg/utl/jwt"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/jpurdie/authapi/pkg/utl/rbac"
	"github.com/jpurdie/authapi/pkg/utl/secure"
	"github.com/jpurdie/authapi/pkg/utl/server"
	"github.com/jpurdie/authapi/pkg/utl/zlog"
	"os"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.Init()
	if err != nil {
		return err
	}
	sec := secure.New(cfg.App.MinPasswordStr, sha1.New())
	rbac := rbac.Service{}
	jwt, err := jwt.New(cfg.JWT.SigningAlgorithm, os.Getenv("JWT_SECRET"), cfg.JWT.DurationMinutes, cfg.JWT.MinSecretLength)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)
	public := e.Group("public")

	pingt.NewHTTP(pingl.New(ping.Initialize(), log), public)

	authMiddleware := authMw.Middleware(jwt)

	//at.NewHTTP(al.New(auth.Initialize(db, jwt, sec, rbac), log), e, authMiddleware)

	v1 := e.Group("api/v1")

	companyTransp.NewHTTP(companyLog.New(company.Initialize(db, rbac, sec), log), v1)

	v1.Use(authMiddleware)
	//ut.NewHTTP(ul.New(user.Initialize(db, rbac, sec), log), v1)

	//	pt.NewHTTP(pl.New(password.Initialize(db, rbac, sec), log), v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
