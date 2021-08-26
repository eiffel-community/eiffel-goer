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
package search

import (
	"net/http"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database"
	"github.com/eiffel-community/eiffel-goer/internal/responses"
)

type SearchHandler struct {
	Config   config.Config
	Database database.Database
}

// Get a new handler for the search endpoint.
func Get(cfg config.Config, db database.Database) *SearchHandler {
	return &SearchHandler{
		cfg, db,
	}
}

// Read handles GET requests against the /search/{id} endpoint.
// To get an event based on eventId passed
func (h *SearchHandler) Read(w http.ResponseWriter, r *http.Request) {
	responses.RespondWithError(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}

// UpstreamDownstream handles POST requests against the /search/{id} endpoint.
// To get upstream/downstream events for an event based on the searchParameters passed
func (h *SearchHandler) UpstreamDownstream(w http.ResponseWriter, r *http.Request) {
	responses.RespondWithError(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}
