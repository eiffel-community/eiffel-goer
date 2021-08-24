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

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database"
	"github.com/eiffel-community/eiffel-goer/internal/responses"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type EventHandler struct {
	Config   config.Config
	Database database.Database
}

// Create a new handler for the event endpoint.
func Get(cfg config.Config, db database.Database) *EventHandler {
	return &EventHandler{
		cfg, db,
	}
}

type EventsSingleRequest struct {
	ID       string `schema:"id"`
	Shallow  bool   `schema:"shallow"`  // TODO: Unused
	Readable bool   `schema:"readable"` // TODO: NYI
}

// Handle GET requests against the /events/{id} endpoint.
// To get single event information
func (h *EventHandler) Read(w http.ResponseWriter, r *http.Request) {
	var request EventsSingleRequest
	if err := schema.NewDecoder().Decode(&request, r.URL.Query()); err != nil {
		responses.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	vars := mux.Vars(r)
	request.ID = vars["id"]
	event, err := h.Database.GetEventByID(request.ID)
	if err != nil {
		responses.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	responses.RespondWithJSON(w, http.StatusOK, event)
}

// Handle GET requests against the /events/
// To get all events information
func (h *EventHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	responses.RespondWithError(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}
