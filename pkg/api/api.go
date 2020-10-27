package api

import (
	"github.com/jpurdie/authapi/pkg/api/invitation"
	invitationl "github.com/jpurdie/authapi/pkg/api/invitation/logging"
	invitationt "github.com/jpurdie/authapi/pkg/api/invitation/transport"
	"github.com/jpurdie/authapi/pkg/api/organization"
	orgl "github.com/jpurdie/authapi/pkg/api/organization/logging"
	orgt "github.com/jpurdie/authapi/pkg/api/organization/transport"
	"github.com/jpurdie/authapi/pkg/api/ping"
	pingl "github.com/jpurdie/authapi/pkg/api/ping/logging"
	pingt "github.com/jpurdie/authapi/pkg/api/ping/transport"
	"github.com/jpurdie/authapi/pkg/api/profile"
	profilel "github.com/jpurdie/authapi/pkg/api/profile/logging"
	profilet "github.com/jpurdie/authapi/pkg/api/profile/transport"
	"github.com/jpurdie/authapi/pkg/api/project"
	projectL "github.com/jpurdie/authapi/pkg/api/project/logging"
	projectT "github.com/jpurdie/authapi/pkg/api/project/transport"
	"github.com/jpurdie/authapi/pkg/api/user"
	userl "github.com/jpurdie/authapi/pkg/api/user/logging"
	"github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"
	usert "github.com/jpurdie/authapi/pkg/api/user/transport"
	"github.com/jpurdie/authapi/pkg/utl/config"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/postgres"
	"github.com/jpurdie/authapi/pkg/utl/server"
	"github.com/jpurdie/authapi/pkg/utl/zlog"

	"os"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	db, err := postgres.New(os.Getenv("DATABASE_URL"), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	e.Static("/swaggerui", cfg.App.SwaggerUIPath)

	//appAPI, err := app.NewAPI(db)
	if err != nil {
		panic(err)
	}

	v1 := e.Group("/api/v1")

	/* Begin Ping Logic */
	pingt.NewHTTP(pingl.New(ping.Initialize(db), log), v1)

	/* Begin User Logic */
	userStruct := user.Initialize(db)
	userLogger := userl.New(userStruct, log)
	usert.NewHTTP(userLogger, v1, db)

	/* Begin Invitation Logic */
	inviteStruct := invitation.Initialize(db, pgsql.User{})
	invitationLogger := invitationl.New(inviteStruct, log)
	invitationt.NewHTTP(invitationLogger, v1, db)

	/* Begin Project Logic */
	projectStruct := project.Initialize(db)
	projectLogger := projectL.New(projectStruct, log)
	projectT.NewHTTP(projectLogger, v1, db)

	orgt.NewHTTP(orgl.New(organization.Initialize(db), log), v1)
	profilet.NewHTTP(profilel.New(profile.Initialize(db), log), v1)

	//everything after here requires auth
	authMiddleware := authMw.Authenticate()
	v1.Use(authMiddleware)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
