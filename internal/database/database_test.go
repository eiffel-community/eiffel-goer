// Copyright 2021 Axis Communications AB.
//
// For a full list of individual contributors, please see the commit history.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package database_test

import (
	"context"
	"testing"

	"github.com/eiffel-community/eiffel-goer/internal/database"
	"github.com/eiffel-community/eiffel-goer/test"
	"github.com/eiffel-community/eiffel-goer/test/mock_drivers"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDriver := mock_drivers.NewMockDatabaseDriver(ctrl)
	mockDB := mock_drivers.NewMockDatabase(ctrl)

	test.SetDatabaseDriver(mockDriver)
	defer test.ResetDatabaseDriver()

	url := "mongodb://db/test"
	entry := &log.Entry{}

	mockDriver.EXPECT().SupportsScheme("mongodb").Return(true)
	mockDriver.EXPECT().Get(gomock.Any(), gomock.Any(), entry).Return(mockDB, nil)
	ctx := context.Background()
	_, err := database.Get(ctx, url, entry)
	assert.NoError(t, err)
}

func TestGetUnknownScheme(t *testing.T) {
	ctx := context.Background()
	_, err := database.Get(ctx, "unknown://db/test", &log.Entry{})
	assert.Error(t, err)
}

func TestGetUnparsableScheme(t *testing.T) {
	ctx := context.Background()
	_, err := database.Get(ctx, "://", &log.Entry{})
	assert.Error(t, err)
}
