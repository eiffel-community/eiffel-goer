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
package api

import (
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/pkg/v1alpha1/handlers/events"
	"github.com/eiffel-community/eiffel-goer/pkg/v1alpha1/handlers/search"
)

type V1Alpha1Application struct {
	Database drivers.Database
	Config   config.Config
	Logger   *log.Entry
}

// Add routes for all handlers to the router.
func (app *V1Alpha1Application) AddRoutes(router *mux.Router) {
	eventHandler := events.Get(app.Config, app.Database, app.Logger)
	searchHandler := search.Get(app.Config, app.Database, app.Logger)

	router.HandleFunc("/events", eventHandler.ReadAll).Methods("GET", "OPTIONS")
	router.HandleFunc("/events/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", eventHandler.Read).Methods("GET", "OPTIONS")
	router.HandleFunc("/search/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", searchHandler.UpstreamDownstream).Methods("POST", "OPTIONS")
}
