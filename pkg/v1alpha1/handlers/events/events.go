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
package events

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/internal/responses"
)

type EventHandler struct {
	Config   config.Config
	Database drivers.Database
	Logger   *log.Entry
}

// Create a new handler for the event endpoint.
func Get(cfg config.Config, db drivers.Database, logger *log.Entry) *EventHandler {
	return &EventHandler{
		cfg, db, logger,
	}
}

type EventsSingleRequest struct {
	Shallow bool `schema:"shallow"` // TODO: Unused
}

// Read handles GET requests against the /events/{id} endpoint.
// To get single event information
func (h *EventHandler) Read(w http.ResponseWriter, r *http.Request) {
	var request EventsSingleRequest
	if err := schema.NewDecoder().Decode(&request, r.URL.Query()); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	vars := mux.Vars(r)
	ID := vars["id"]
	event, err := h.Database.GetEventByID(r.Context(), ID)
	if err != nil {
		responses.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	responses.RespondWithJSON(w, http.StatusOK, event)
}

// ReadAll handles GET requests against the /events/ endpoint.
// To get all events information
func (h *EventHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	responses.RespondWithError(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}
