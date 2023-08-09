// Copyright 2021 Axis Communications AB.
//
// For a full list of individual contributors, please see the commit history.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package application

import (
	"context"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/pkg/server"
	v1api "github.com/eiffel-community/eiffel-goer/pkg/v1/api"
)

type Application struct {
	Database drivers.Database
	Config   config.Config
	Router   *mux.Router
	Server   server.Server
	V1       *v1api.V1Application
	Logger   *log.Entry
}

// Get a new Goer application.
func Get(ctx context.Context, cfg config.Config, logger *log.Entry) (*Application, error) {
	application := &Application{
		Config: cfg,
		Router: mux.NewRouter(),
		Server: server.Get(),
		Logger: logger,
	}
	if cfg.DBConnectionString() != "" {
		db, err := application.getDB(ctx)
		if err != nil {
			return nil, err
		}
		application.Database = db
	}
	return application, nil
}

// getDB gets, and connects to, a database.
func (app *Application) getDB(ctx context.Context) (drivers.Database, error) {
	db, err := database.Get(
		ctx,
		app.Config.DBConnectionString(),
		app.Logger,
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// LoadV1Routes loads routes for the /v1/ endpoint.
func (app *Application) LoadV1Routes() {
	app.V1 = &v1api.V1Application{
		Config:   app.Config,
		Database: app.Database,
		Logger:   app.Logger,
	}
	subrouter := app.Router.PathPrefix("/v1").Name("v1").Subrouter()
	app.V1.AddRoutes(subrouter)
}

// Start connects to the database and starts the webserver.
// This is a blocking function, waiting for the webserver to shut down.
func (app *Application) Start(ctx context.Context) error {
	srv := app.Server.WithAddr(app.Config.APIPort()).WithRouter(app.Router)
	defer func() {
		if err := app.Stop(ctx); err != nil {
			app.Logger.Errorf("Error stopping application: %s", err)
		}
	}()
	if err := srv.Start(); err != nil {
		return err
	}
	srv.WaitStopped()
	return srv.Error()
}

// Stop the application and close the database connection.
func (app *Application) Stop(ctx context.Context) error {
	return app.Database.Close(ctx)
}
