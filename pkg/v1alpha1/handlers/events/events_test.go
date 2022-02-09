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
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eiffel-community/eiffelevents-sdk-go"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/test/mock_config"
	"github.com/eiffel-community/eiffel-goer/test/mock_drivers"
)

var activityJSON = []byte(`
{
    "data": {
        "name": "Test activity"
    },
    "links": [],
    "meta": {
        "id": "e04cf9d3-4d57-471e-bd65-f8fc20d21d84",
        "time": 1629449650361,
        "type": "EiffelActivityTriggeredEvent",
        "version": "3.0.0"
    }
}
`)

// Test that the events/{id} endpoint work as expected.
func TestEvents(t *testing.T) {
	// Load the test event twice; once as a drivers.EiffelEvent and once as
	// a proper struct via eiffelevents. The latter is only used to easily
	// extract the event ID.
	eventMap := make(drivers.EiffelEvent)
	require.NoError(t, json.Unmarshal(activityJSON, &eventMap))
	event, err := eiffelevents.UnmarshalAny(activityJSON)
	require.NoError(t, err)
	eventID := event.(eiffelevents.MetaTeller).ID()

	badRequest := httptest.NewRequest(http.MethodGet, "/events/"+eventID, nil)
	q := badRequest.URL.Query()
	q.Add("nah", "hello")
	badRequest.URL.RawQuery = q.Encode()

	tests := []struct {
		name       string
		request    *http.Request
		statusCode int
		eventID    string
		mockError  error
	}{
		{name: "Read", request: httptest.NewRequest(http.MethodGet, "/events/"+eventID, nil), statusCode: http.StatusOK, eventID: eventID},
		{name: "ReadBadRequest", request: badRequest, statusCode: http.StatusBadRequest, eventID: ""},
		{name: "ReadNotFound", request: httptest.NewRequest(http.MethodGet, "/events/"+eventID, nil), statusCode: http.StatusNotFound, eventID: "", mockError: errors.New("not found")},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockCfg := mock_config.NewMockConfig(ctrl)
			mockDB := mock_drivers.NewMockDatabase(ctrl)
			mockDB.EXPECT().GetEventByID(gomock.Any(), eventID).Return(eventMap, testCase.mockError)
			app := Get(mockCfg, mockDB, &log.Entry{})
			handler := mux.NewRouter()
			handler.HandleFunc("/events/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", app.Read)

			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, testCase.request)

			assert.Equalf(t, testCase.statusCode, responseRecorder.Code, "Input URL: %s", testCase.request.URL)
			if responseRecorder.Code == http.StatusOK {
				assert.JSONEq(t, string(activityJSON), responseRecorder.Body.String())
			}
		})
	}
}
