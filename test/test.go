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

// This package is for generating the mocks that are used in test as well
// as global test helper functions.
package test

import (
	"github.com/eiffel-community/eiffel-goer/internal/database"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
)

//go:generate mockgen -destination mock_drivers/mock_drivers.go github.com/eiffel-community/eiffel-goer/internal/database/drivers DatabaseDriver,Database
//go:generate mockgen -destination mock_config/mock_config.go github.com/eiffel-community/eiffel-goer/internal/config Config
//go:generate mockgen -destination mock_server/mock_server.go github.com/eiffel-community/eiffel-goer/pkg/server Server

var copiedDrivers []drivers.DatabaseDriver

// SetDatabaseDriver sets the slice of available database drivers to a mocked
// version of the database driver. This call shall be followed by a defer call
// to ResetDatabaseDriver or all subsequent tests which use the database will fail.
func SetDatabaseDriver(mockDriver drivers.DatabaseDriver) {
	copiedDrivers = database.Drivers
	database.Drivers = []drivers.DatabaseDriver{mockDriver}
}

// ResetDatabaseDriver will set the available database drivers back to the original
// value that was set when SetDatabaseDriver was called.
// This function should be defer called after SetDatabaseDriver was called.
func ResetDatabaseDriver() {
	database.Drivers = copiedDrivers
	copiedDrivers = []drivers.DatabaseDriver{}
}
