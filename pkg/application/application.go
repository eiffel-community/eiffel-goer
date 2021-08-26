// Copyright 2021 Axis Communications AB.
//
// For a full list of individual contributors, please see the commit history.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database"
	"github.com/eiffel-community/eiffel-goer/pkg/server"
	v1alpha1 "github.com/eiffel-community/eiffel-goer/pkg/v1alpha1/api"
)

type Application struct {
	Database database.Database
	Config   config.Config
	Router   *mux.Router
	Server   server.Server
	V1Alpha1 *v1alpha1.V1Alpha1Application
}

// Get a new Goer application.
func Get(cfg config.Config) (*Application, error) {
	application := &Application{
		Config: cfg,
		Router: mux.NewRouter(),
		Server: server.Get(),
	}
	if cfg.DBConnectionString() != "" {
		db, err := application.getDB()
		if err != nil {
			return nil, err
		}
		application.Database = db
	}
	return application, nil
}

// getDB gets, but does not connect to, a database.
func (app *Application) getDB() (database.Database, error) {
	db, err := database.Get(
		app.Config.DBConnectionString(),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// LoadV1Alpha1Routes loads routes for the /v1alpha1/ endpoint.
func (app *Application) LoadV1Alpha1Routes() {
	app.V1Alpha1 = &v1alpha1.V1Alpha1Application{
		Config:   app.Config,
		Database: app.Database,
	}
	subrouter := app.Router.PathPrefix("/v1alpha1").Name("v1alpha1").Subrouter()
	app.V1Alpha1.AddRoutes(subrouter)
}

// Start connects to the database and starts the webserver.
// This is a blocking function, waiting for the webserver to shut down.
func (app *Application) Start(ctx context.Context) error {
	srv := app.Server.WithAddr(app.Config.APIPort()).WithRouter(app.Router)
	err := app.Database.Connect(ctx)
	if err != nil {
		return err
	}
	defer app.Stop(ctx)
	err = srv.Start()
	if err != nil {
		return err
	}
	srv.WaitStopped()
	return srv.Error()
}

// Stop the application and close the database connection.
func (app *Application) Stop(ctx context.Context) error {
	return app.Database.Close(ctx)
}
