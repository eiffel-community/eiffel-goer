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

// This package tests that all endpoints are added and respond to requests.
// The function of each handler is tested separately in each handler.
package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/schema"
	"github.com/eiffel-community/eiffel-goer/pkg/application"
	"github.com/eiffel-community/eiffel-goer/test/mock_config"
	"github.com/eiffel-community/eiffel-goer/test/mock_drivers"
)

// Test that all v1alpha1 endpoints are added properly.
func TestRoutes(t *testing.T) {
	eventID := "3fabaa6b-5343-4d74-8af9-dc2e4c1f2827"
	tests := []struct {
		name       string
		url        string
		httpMethod string
		statusCode int
	}{
		{name: "EventsRead", httpMethod: http.MethodGet, url: "/v1alpha1/events/" + eventID, statusCode: http.StatusOK},
		{name: "EventsReadAll", httpMethod: http.MethodGet, url: "/v1alpha1/events", statusCode: http.StatusNotImplemented},
		{name: "SearchRead", httpMethod: http.MethodGet, url: "/v1alpha1/search/" + eventID, statusCode: http.StatusNotImplemented},
		{name: "SearchUpstreamDownstream", httpMethod: http.MethodPost, url: "/v1alpha1/search/" + eventID, statusCode: http.StatusNotImplemented},
	}

	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabaseDriver(ctrl)

	mockCfg.EXPECT().DBConnectionString().Return("").AnyTimes()
	mockCfg.EXPECT().APIPort().Return(":8080").AnyTimes()
	// Have to use 'gomock.Any()' for the context as mux adds values to the request context.
	mockDB.EXPECT().GetEventByID(gomock.Any(), eventID).Return(schema.EiffelEvent{}, nil)
	mockDB.EXPECT().GetEvents(gomock.Any()).Return([]schema.EiffelEvent{}, nil)
	mockDB.EXPECT().UpstreamDownstreamSearch(gomock.Any(), "id").Return([]schema.EiffelEvent{}, nil)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			app, err := application.Get(mockCfg, &log.Entry{})
			if err != nil {
				t.Error(err)
			}
			app.Database = mockDB
			app.LoadV1Alpha1Routes()

			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(testCase.httpMethod, testCase.url, nil)

			app.Router.ServeHTTP(responseRecorder, request)
			expectedStatusCode := testCase.statusCode
			if responseRecorder.Code != expectedStatusCode {
				t.Errorf("Want status '%d' for %q, got '%d'", expectedStatusCode, testCase.url, responseRecorder.Code)
			}

		})
	}
}
