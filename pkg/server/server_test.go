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
	"log"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert.IsType(t, &WebServer{}, Get())
}

func TestWithAddr(t *testing.T) {
	addr := "127.0.0.1:8080"
	server, ok := Get().WithAddr(addr).(*WebServer)
	assert.Truef(t, ok, "WebServer struct is not a Server interface")
	assert.Equal(t, addr, server.server.Addr)
}

func TestWithErrLogger(t *testing.T) {
	logger := &log.Logger{}
	server, ok := Get().WithErrLogger(logger).(*WebServer)
	assert.Truef(t, ok, "WebServer struct is not a Server interface")
	assert.Equal(t, logger, server.server.ErrorLog)
}

func TestWithRouter(t *testing.T) {
	router := mux.NewRouter()
	server, ok := Get().WithRouter(router).(*WebServer)
	assert.Truef(t, ok, "WebServer struct is not a Server interface")
	assert.Equal(t, router, server.server.Handler)
}

func TestStart(t *testing.T) {
	server := Get().WithAddr("127.0.0.1:8080").WithRouter(mux.NewRouter())
	err := server.Start()
	assert.NoError(t, err)
	assert.Truef(t, server.WaitRunning(), "server did not start properly")
	server.Close()
	assert.Truef(t, server.WaitStopped(), "server did not stop properly")
	assert.ErrorIs(t, server.Error(), http.ErrServerClosed)
}

func TestStartNoAddr(t *testing.T) {
	server := Get().WithRouter(mux.NewRouter())
	err := server.Start()
	assert.Error(t, err)
	assert.Equal(t, "server missing address", err.Error())
	assert.False(t, server.WaitRunning())
}

func TestStartNoHandler(t *testing.T) {
	server := Get().WithAddr("127.0.0.1:8080")
	err := server.Start()
	assert.Error(t, err)
	assert.Equal(t, "server missing handler", err.Error())
	assert.False(t, server.WaitRunning())
}
