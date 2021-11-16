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
	"github.com/stretchr/testify/assert"

	"github.com/eiffel-community/eiffel-goer/test"
	"github.com/eiffel-community/eiffel-goer/test/mock_config"
	"github.com/eiffel-community/eiffel-goer/test/mock_drivers"
	"github.com/eiffel-community/eiffel-goer/test/mock_server"
)

// Test that it is possible to get an application.
func TestGet(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)

	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDB, nil)
	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	app, err := Get(ctx, mockCfg, &log.Entry{})
	assert.NoError(t, err)
	assert.Equal(t, mockCfg, app.Config)
	assert.NotNil(t, app.Database)
	assert.NotNil(t, app.Router)
	assert.NotNil(t, app.Server)
}

// Test that it is possible to get an application without a database.
func TestGetNoDB(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("")

	app, err := Get(ctx, mockCfg, &log.Entry{})
	assert.NoError(t, err)
	assert.Equal(t, mockCfg, app.Config)
	assert.Nil(t, app.Database)
	assert.NotNil(t, app.Router)
	assert.NotNil(t, app.Server)
}

// Test that Get return error if there was an error when getting database.
func TestGetDBError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("invalid://testdb").Times(2)

	_, err := Get(ctx, mockCfg, &log.Entry{})
	assert.Errorf(t, err, "application should have raised error due to invalid database connection string")
}

// Test that getDB return a database interface.
func TestGetDB(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)

	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDB, nil)
	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	application := &Application{
		Config: mockCfg,
	}

	_, err := application.getDB(ctx)
	assert.NoError(t, err)
}

// Test that the application creates the v1alpha1 subrouter.
func TestLoadV1Alpha1Routes(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)

	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDB, nil)
	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	app, err := Get(ctx, mockCfg, &log.Entry{})
	assert.NoError(t, err)

	app.LoadV1Alpha1Routes()
	assert.NotNil(t, app.Router.Get("v1alpha1"))
}

// Test that the application starts the WebServer & connects to the Database.
func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	mockServer := mock_server.NewMockServer(ctrl)
	ctx := context.Background()

	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDB, nil)
	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	mockDB.EXPECT().Close(gomock.Any()).Return(nil)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockCfg.EXPECT().APIPort().Return(":8080")

	app, err := Get(ctx, mockCfg, &log.Entry{})
	assert.NoError(t, err)

	mockServer.EXPECT().WithAddr(":8080").Return(mockServer)
	mockServer.EXPECT().WithRouter(app.Router).Return(mockServer)
	mockServer.EXPECT().Start().Return(nil)
	mockServer.EXPECT().WaitStopped().Return(true)
	mockServer.EXPECT().Error().Return(nil)

	app.Server = mockServer

	assert.NoError(t, app.Start(ctx))
}

// Test that application returns error if server start fails.
func TestStartFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	mockServer := mock_server.NewMockServer(ctrl)
	ctx := context.Background()

	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDB, nil)
	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	mockDB.EXPECT().Close(gomock.Any()).Return(nil)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)
	mockCfg.EXPECT().APIPort().Return("")

	app, err := Get(ctx, mockCfg, &log.Entry{})
	assert.NoError(t, err)

	mockServer.EXPECT().WithAddr("").Return(mockServer)
	mockServer.EXPECT().WithRouter(app.Router).Return(mockServer)
	mockServer.EXPECT().Start().Return(errors.New("error starting"))

	app.Server = mockServer

	assert.Error(t, app.Start(ctx))
}

// Test that application closes the database on Stop.
func TestStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCfg := mock_config.NewMockConfig(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)
	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	ctx := context.Background()

	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDB, nil)
	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	mockDB.EXPECT().Close(gomock.Any()).Return(nil)
	mockCfg.EXPECT().DBConnectionString().Return("mongodb://testdb/testdb").Times(2)

	app, err := Get(ctx, mockCfg, &log.Entry{})
	assert.NoError(t, err)

	assert.NoError(t, app.Stop(ctx))
}
