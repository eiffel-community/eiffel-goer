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
package application

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/test/mock_config"
	"github.com/eiffel-community/eiffel-goer/test/mock_drivers"
	"github.com/eiffel-community/eiffel-goer/test/mock_server"
)

// Test that it is possible to get an application.
func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)

	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}
	if app.Config != mockCfg {
		t.Errorf("config not set properly by Get")
	}
	if app.Database == nil {
		t.Error("application did not set up database")
	}
	if app.Router == nil {
		t.Error("application did not set up router")
	}
	if app.Server == nil {
		t.Error("application did not set up server")
	}
}

// Test that it is possible to get an application without a database.
func TestGetNoDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("")

	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}
	if app.Config != mockCfg {
		t.Errorf("config not set properly by Get")
	}
	if app.Database != nil {
		t.Error("application did not set up database")
	}
	if app.Router == nil {
		t.Error("application did not set up router")
	}
	if app.Server == nil {
		t.Error("application did not set up server")
	}
}

// Test that Get return error if there was an error when getting database.
func TestGetDBError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("invalid://testdb").Times(2)

	_, err := Get(mockCfg, &log.Entry{})
	if err == nil {
		t.Error("application should have raised error due to invalid database connection string")
	}
}

// Test that getDB return a database interface.
func TestGetDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)

	application := &Application{
		Config: mockCfg,
	}
	db, err := application.getDB()
	if err != nil {
		t.Error(err)
	}
	_, ok := db.(drivers.Database)
	if !ok {
		t.Error("database from 'getDB' is not a Database interface")
	}
}

// Test that the application creates the v1alpha1 subrouter.
func TestLoadV1Alpha1Routes(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}

	err = app.LoadV1Alpha1Routes()
	if err != nil {
		t.Error(err)
	}
	route := app.Router.Get("v1alpha1")
	if route == nil {
		t.Error("the v1alpha1 route did not get loaded")
	}
}

// Test that the application exits if there is an error getting the database driver.
func TestLoadV1Alpha1RoutesDriverError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockDB.EXPECT().Driver().Return(&mock_drivers.MockDatabaseDriver{}, errors.New("this is an error"))
	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}
	app.Database = mockDB

	err = app.LoadV1Alpha1Routes()
	if err == nil {
		t.Error("the v1alpha1 routes loaded when there was an error getting database driver")
	}
	route := app.Router.Get("v1alpha1")
	if route != nil {
		t.Error("the v1alpha1 route did not get loaded")
	}
}

// Test that the application starts the WebServer & connects to the Database.
func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockServer := mock_server.NewMockServer(ctrl)
	ctx := context.Background()

	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockCfg.EXPECT().APIPort().Return(":8080")

	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}

	mockDB.EXPECT().Connect(ctx).Return(nil)
	mockDB.EXPECT().Close(ctx).Return(nil)
	mockServer.EXPECT().WithAddr(":8080").Return(mockServer)
	mockServer.EXPECT().WithRouter(app.Router).Return(mockServer)
	mockServer.EXPECT().Start().Return(nil)
	mockServer.EXPECT().WaitStopped().Return(true)
	mockServer.EXPECT().Error().Return(nil)

	app.Database = mockDB
	app.Server = mockServer

	err = app.Start(ctx)
	if err != nil {
		t.Error(err)
	}
}

// Test that the application Start aborts with error if database connect fails
func TestStartAbort(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockServer := mock_server.NewMockServer(ctrl)
	ctx := context.Background()

	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockCfg.EXPECT().APIPort().Return(":8080")

	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}

	mockServer.EXPECT().WithAddr(":8080").Return(mockServer)
	mockServer.EXPECT().WithRouter(app.Router).Return(mockServer)
	mockServer.EXPECT().WaitStopped().Return(true)
	mockServer.EXPECT().Error().Return(nil)
	mockDB.EXPECT().Connect(ctx).Return(errors.New("did not work"))

	app.Database = mockDB
	app.Server = mockServer

	err = app.Start(ctx)
	if err == nil {
		t.Error("application did not abort start after error on database.Connect")
	}
}

// Test that application returns error if server start fails.
func TestStartFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockServer := mock_server.NewMockServer(ctrl)
	ctx := context.Background()

	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockCfg.EXPECT().APIPort().Return("")

	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}

	mockDB.EXPECT().Connect(ctx).Return(nil)
	mockDB.EXPECT().Close(ctx).Return(nil)
	mockServer.EXPECT().WithAddr("").Return(mockServer)
	mockServer.EXPECT().WithRouter(app.Router).Return(mockServer)
	mockServer.EXPECT().Start().Return(errors.New("error starting"))

	app.Database = mockDB
	app.Server = mockServer

	err = app.Start(ctx)
	if err == nil {
		t.Error("application start did not abort when server.Start failed")
	}
}

// Test that application closes the database on Stop.
func TestStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	ctx := context.Background()

	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)

	app, err := Get(mockCfg, &log.Entry{})
	if err != nil {
		t.Error(err)
	}

	mockDB.EXPECT().Close(ctx).Return(nil)
	app.Database = mockDB
	app.Stop(ctx)
}
