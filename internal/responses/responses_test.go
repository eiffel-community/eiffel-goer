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
package responses

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// Test that RespondWithJSON writes the correct HTTP code, message and adds a content type header.
func TestRespondWithJSON(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	RespondWithJSON(responseRecorder, 200, map[string]string{"hello": "world"})
	if responseRecorder.Header().Get("Content-Type") != "application/json" {
		t.Error("content type is not set to application/json")
	}
	if responseRecorder.Result().StatusCode != 200 {
		t.Error("status code was not properly set")
	}
	var jsonBody map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &jsonBody)
	if err != nil {
		t.Error(err)
	}
	if jsonBody["hello"] != "world" {
		t.Error("body was set incorrectly")
	}
}

// Test that RespondWithError writes the correct HTTP code, message and adds a content type header.
func TestRespondWithError(t *testing.T) {
	responseRecorder := httptest.NewRecorder()
	RespondWithError(responseRecorder, 400, "failure")
	if responseRecorder.Header().Get("Content-Type") != "application/json" {
		t.Error("content type is not set to application/json")
	}
	if responseRecorder.Result().StatusCode != 400 {
		t.Error("status code was not properly set")
	}
	var jsonBody map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &jsonBody)
	if err != nil {
		t.Error(err)
	}
	if jsonBody["error"] != "failure" {
		t.Error("body was set incorrectly")
	}
}
