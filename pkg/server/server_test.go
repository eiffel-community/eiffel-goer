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
package server

import (
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestGet(t *testing.T) {
	_, ok := Get().(*WebServer)
	if !ok {
		t.Error("WebServer struct is not a Server interface")
	}
}

func TestWithAddr(t *testing.T) {
	addr := "127.0.0.1:8080"
	server, ok := Get().WithAddr(addr).(*WebServer)
	if !ok {
		t.Error("WebServer struct is not a Server interface")
	}
	if server.server.Addr != addr {
		t.Error("address was not registered on server")
	}
}

func TestWithErrLogger(t *testing.T) {
	logger := &log.Logger{}
	server, ok := Get().WithErrLogger(logger).(*WebServer)
	if !ok {
		t.Error("WebServer struct is not a Server interface")
	}
	if server.server.ErrorLog != logger {
		t.Error("logger was not registered on server")
	}
}

func TestWithRouter(t *testing.T) {
	router := mux.NewRouter()
	server, ok := Get().WithRouter(router).(*WebServer)
	if !ok {
		t.Error("WebServer struct is not a Server interface")
	}
	if server.server.Handler != router {
		t.Error("router was not registered on server")
	}
}

func TestStart(t *testing.T) {
	server := Get().WithAddr("127.0.0.1:8080").WithRouter(mux.NewRouter())
	err := server.Start()
	if err != nil {
		t.Error(err)
	}
	if !server.WaitRunning() {
		t.Error("server did not start properly")
	}
	server.Close()
	if !server.WaitStopped() {
		t.Error("server did not stop properly")
	}
	err = server.Error()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Errorf("there was an error in the server %s", err)
	}
}

func TestStartNoAddr(t *testing.T) {
	server := Get().WithRouter(mux.NewRouter())
	err := server.Start()
	if err == nil {
		t.Error(err)
	}
	if err.Error() != "server missing address" {
		t.Errorf("error message does not specify it's missing address: %q", err.Error())
	}
	if server.WaitRunning() {
		t.Error("server started when missing addr")
	}
}

func TestStartNoHandler(t *testing.T) {
	server := Get().WithAddr("127.0.0.1:8080")
	err := server.Start()
	if err == nil {
		t.Error(err)
	}
	if err.Error() != "server missing handler" {
		t.Errorf("error message does not specify it's missing handler: %q", err.Error())
	}
	if server.WaitRunning() {
		t.Error("server started when missing handler")
	}
}
