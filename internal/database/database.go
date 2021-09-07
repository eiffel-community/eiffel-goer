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
package database

import (
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/internal/database/drivers/mongodb"
)

// Get a new Database.
func Get(connectionString string, logger *log.Entry) (drivers.Database, error) {
	var databases = []drivers.Database{&mongodb.Database{}}

	connectionURL, err := url.Parse(connectionString)
	if err != nil {
		return nil, err
	}
	for _, database := range databases {
		if database.SupportsScheme(connectionURL.Scheme) {
			return database.Get(connectionURL, logger)
		}
	}
	return nil, fmt.Errorf("cannot find database for scheme %q", connectionURL.Scheme)
}
