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

// This package implements the database interface against MongoDB following
// the collection structure implemented by the Eiffel GraphQL API and
// Simple Event Sender.
package mongodb

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/pkg/schema"
)

type Driver struct {
	Client   *mongo.Client
	Database *mongo.Database
	Logger   *log.Entry
}

// Get creates a new database.Database interface against MongoDB.
func (m *Driver) Get(connectionURL *url.URL, logger *log.Entry) (drivers.DatabaseDriver, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURL.String()))
	if err != nil {
		return nil, err
	}
	// The Path value from url.URL always has a '/' prepended.
	databaseName := strings.Split(connectionURL.Path, "/")[1]
	return &Driver{
		Client:   client,
		Database: client.Database(databaseName),
		Logger:   logger,
	}, nil
}

// Test whether the MongoDB driver supports a scheme.
func (m *Driver) SupportsScheme(scheme string) bool {
	switch scheme {
	case "mongodb":
		return true
	case "mongodb+srv":
		return true
	default:
		return false
	}
}

// Connect to the MongoDB database and ping it to make sure it works.
func (m *Driver) Connect(ctx context.Context) error {
	err := m.Client.Connect(ctx)
	if err != nil {
		return err
	}
	return m.Client.Ping(ctx, readpref.Primary())
}

// GetEvents gets all events information.
func (m *Driver) GetEvents(ctx context.Context) ([]schema.EiffelEvent, error) {
	return nil, errors.New("not yet implemented")
}

// SearchEvent searches for an event based on event ID.
func (m *Driver) SearchEvent(ctx context.Context, id string) (schema.EiffelEvent, error) {
	return schema.EiffelEvent{}, errors.New("not yet implemented")
}

// UpstreamDownstreamSearch searches for events upstream and/or downstream of event by ID.
func (m *Driver) UpstreamDownstreamSearch(ctx context.Context, id string) ([]schema.EiffelEvent, error) {
	return nil, errors.New("not yet implemented")
}

// GetEventByID gets an event by ID in all collections.
func (m *Driver) GetEventByID(ctx context.Context, id string) (schema.EiffelEvent, error) {
	collections, err := m.Database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return schema.EiffelEvent{}, err
	}
	filter := bson.D{{"meta.id", id}}
	for _, collection := range collections {
		var event schema.EiffelEvent
		singleResult := m.Database.Collection(collection).FindOne(ctx, filter)
		err := singleResult.Decode(&event)
		if err != nil {
			continue
		} else {
			return event, nil
		}
	}
	return schema.EiffelEvent{}, fmt.Errorf("%q not found in any collection", id)
}

// Close the database connection.
func (m *Driver) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
