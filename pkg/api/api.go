package api

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jpurdie/authapi/pkg/api/invitation"
	invitationL "github.com/jpurdie/authapi/pkg/api/invitation/logging"
	invitationT "github.com/jpurdie/authapi/pkg/api/invitation/transport"
	userDB "github.com/jpurdie/authapi/pkg/api/user/platform/pgsql"

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
	usert "github.com/jpurdie/authapi/pkg/api/user/transport"

	"github.com/jpurdie/authapi/pkg/utl/config"
	authMw "github.com/jpurdie/authapi/pkg/utl/middleware/auth"
	"github.com/jpurdie/authapi/pkg/utl/server"
	"github.com/jpurdie/authapi/pkg/utl/zlog"
	_ "github.com/lib/pq" // here

	"os"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {
	//db, err := postgres.New(os.Getenv("DATABASE_URL"), cfg.DB.Timeout, cfg.DB.LogQueries)
	//if err != nil {
	//	return err
	//}
	dbx, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Print(err)
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
	pingt.NewHTTP(pingl.New(ping.Initialize(dbx), log), v1)

	/* Begin User Logic */
	userStruct := user.Initialize(dbx)
	userLogger := userl.New(userStruct, log)
	usert.NewHTTP(userLogger, v1, dbx)

	///* Begin Invitation Logic */
	inviteStruct := invitation.Initialize(dbx, userDB.User{})
	invitationLogger := invitationL.New(inviteStruct, log)
	invitationT.NewHTTP(invitationLogger, v1, dbx)

	///* Begin Project Logic */
	projectStruct := project.Initialize(dbx)
	projectLogger := projectL.New(projectStruct, log)
	projectT.NewHTTP(projectLogger, v1, dbx)
	//
	///* Begin Strategic Alignment Logic */
	//alignmentStruct := strategicalignment.Initialize(dbx)
	//alignmentLogger := strategicL.New(alignmentStruct, log)
	//strategicT.NewHTTP(alignmentLogger, v1, dbx)

	orgt.NewHTTP(orgl.New(organization.Initialize(dbx), log), v1)
	profilet.NewHTTP(profilel.New(profile.Initialize(dbx), log), v1)

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
