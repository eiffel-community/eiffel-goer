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
package search

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/internal/responses"
)

type Handler struct {
	Config   config.Config
	Database drivers.Database
	Logger   *log.Entry
}

// Get a new handler for the search endpoint.
func Get(cfg config.Config, db drivers.Database, logger *log.Entry) *Handler {
	return &Handler{
		cfg, db, logger,
	}
}

// UpstreamDownstream handles POST requests against the /search/{id} endpoint.
// To get upstream/downstream events for an event based on the searchParameters passed.
func (h *Handler) UpstreamDownstream(w http.ResponseWriter, _ *http.Request) {
	responses.RespondWithError(w, http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))
}
