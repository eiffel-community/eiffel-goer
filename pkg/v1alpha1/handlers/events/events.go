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
	"fmt"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/internal/query"
	"github.com/eiffel-community/eiffel-goer/internal/requests"
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

// Read handles GET requests against the /events/{id} endpoint.
// To get single event information.
func (h *EventHandler) Read(w http.ResponseWriter, r *http.Request) {
	var request requests.SingleEventRequest
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

func getTags(tagName string, item interface{}) map[string]struct{} {
	t := reflect.TypeOf(item).Elem()
	tags := make(map[string]struct{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		tags[tag] = struct{}{}
	}
	return tags
}

type multiResponse struct {
	PageNo           int32                 `json:"pageNo"`
	PageSize         int32                 `json:"pageSize"`
	TotalNumberItems int                   `json:"totalNumberItems"`
	Items            []drivers.EiffelEvent `json:"items"`
}

// buildConditions takes a raw URL query, parses out all conditions and removes ignoreKeys.
func buildConditions(rawQuery string, ignoreKeys map[string]struct{}) ([]query.Condition, error) {
	res, err := query.Parse("nofile", []byte(rawQuery))
	if err != nil {
		return nil, err
	}
	allConditions, ok := res.([]query.Condition)
	if !ok {
		return nil, fmt.Errorf("query parser unexpectedly returned a %T value from the query %q", res, rawQuery)
	}
	var conditions []query.Condition
	for _, condition := range allConditions {
		_, ok := ignoreKeys[condition.Field]
		if !ok {
			conditions = append(conditions, condition)
		}
	}
	return conditions, nil
}

// ReadAll handles GET requests against the /events/ endpoint.
// To get all events information.
func (h *EventHandler) ReadAll(w http.ResponseWriter, r *http.Request) {
	request := requests.MultipleEventsRequest{
		Shallow:       false,
		PageNo:        1,
		PageSize:      500,
		PageStartItem: 1,
		Lazy:          false,
		Readable:      false,
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	if err := decoder.Decode(&request, r.URL.Query()); err != nil {
		h.Logger.Error(err)
		responses.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ignoreKeys := getTags("schema", &request)
	conditions, err := buildConditions(r.URL.RawQuery, ignoreKeys)
	if err != nil {
		h.Logger.Error(err)
		responses.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	request.Conditions = conditions
	events, err := h.Database.GetEvents(r.Context(), request)
	if err != nil {
		h.Logger.Error(err)
		responses.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	response := multiResponse{
		request.PageNo,
		request.PageSize,
		len(events), // TODO: This is not correct at the moment.
		events,
	}
	responses.RespondWithJSON(w, http.StatusOK, response)
}
