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

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/internal/schema"
)

// Database is a MongoDB database connection.
type Database struct {
	logger           *log.Entry
	client           *mongo.Client
	connectionString connstring.ConnString
}

// Get creates a new database.Database interface against MongoDB.
func (d *Database) Get(connectionURL *url.URL, logger *log.Entry) (drivers.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURL.String()))
	if err != nil {
		return nil, err
	}
	connectionString, err := connstring.Parse(connectionURL.String())
	if err != nil {
		return nil, err
	}
	return &Database{
		client:           client,
		connectionString: connectionString,
		logger:           logger,
	}, nil
}

// Test whether the MongoDB driver supports a scheme.
func (d *Database) SupportsScheme(scheme string) bool {
	switch scheme {
	case "mongodb":
		return true
	case "mongodb+srv":
		return true
	default:
		return false
	}
}

// Driver creates a new database driver for event requests.
func (d *Database) Driver() (drivers.DatabaseDriver, error) {
	return &Driver{
		database: d.client.Database(d.connectionString.Database),
		logger:   d.logger,
	}, nil
}

// Connect to the MongoDB database and ping it to make sure it works.
func (d *Database) Connect(ctx context.Context) error {
	err := d.client.Connect(ctx)
	if err != nil {
		return err
	}
	return d.client.Ping(ctx, readpref.Primary())
}

// Close the database connection.
func (d *Database) Close(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}

// Driver is a database driver for requesting events from MongoDB.
type Driver struct {
	database *mongo.Database
	logger   *log.Entry
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
	collections, err := m.database.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return schema.EiffelEvent{}, err
	}
	filter := bson.D{{"meta.id", id}}
	for _, collection := range collections {
		var event schema.EiffelEvent
		singleResult := m.database.Collection(collection).FindOne(ctx, filter)
		err := singleResult.Decode(&event)
		if err != nil {
			continue
		} else {
			return event, nil
		}
	}
	return schema.EiffelEvent{}, fmt.Errorf("%q not found in any collection", id)
}
