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
	"strconv"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	"github.com/eiffel-community/eiffel-goer/internal/database/drivers"
	"github.com/eiffel-community/eiffel-goer/internal/query"
	"github.com/eiffel-community/eiffel-goer/internal/requests"
)

// Database is a MongoDB database connection.
type Driver struct {
	client           *mongo.Client
	connectionString connstring.ConnString
}

// Get creates and connects a new database.Database interface against MongoDB.
func (d *Driver) Get(ctx context.Context, connectionURL *url.URL, logger *log.Entry) (drivers.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionURL.String()))
	if err != nil {
		return nil, err
	}
	d.client = client
	connectionString, err := connstring.Parse(connectionURL.String())
	if err != nil {
		return nil, err
	}
	d.connectionString = connectionString
	return d.connect(ctx, logger)
}

// Test whether the MongoDB driver supports a scheme.
func (d *Driver) SupportsScheme(scheme string) bool {
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
func (d *Driver) connect(ctx context.Context, logger *log.Entry) (drivers.Database, error) {
	err := d.client.Connect(ctx)
	if err != nil {
		return &Database{}, err
	}
	if err = d.client.Ping(ctx, readpref.Primary()); err != nil {
		return &Database{}, err
	}
	return &Database{
		database: d.client.Database(d.connectionString.Database),
		client:   d.client,
		logger:   logger,
	}, nil
}

// Database is a connected database interface for requesting events from MongoDB.
type Database struct {
	database *mongo.Database
	client   *mongo.Client
	logger   *log.Entry
}

// operators is a translation table from query.Param to mongodb operators.
var operators = map[string]string{
	"=": "$eq", "!=": "$ne", ">": "$gt", "<": "$lt", "<=": "$lte", ">=": "$gte", "exists": "$exists",
}

// typeCast values in condition based on TypeConv parameter. Returns a bson Element.
func typeCast(condition query.Condition) (bson.E, error) {
	var err error
	e := bson.E{Key: operators[condition.Op]}
	switch condition.TypeConv {
	case "int":
		e.Value, err = strconv.ParseInt(condition.Value, 0, 64)
	case "double":
		e.Value, err = strconv.ParseFloat(condition.Value, 64)
	case "bool":
		e.Value, err = strconv.ParseBool(condition.Value)
	default:
		e.Value = condition.Value
		err = nil
	}
	return e, err
}

// buildFilter creates a MongoDB filter based on query parameters.
func buildFilter(conditions []query.Condition) (bson.D, error) {
	d := bson.D{}
	elements := map[string]bson.D{}
	for _, condition := range conditions {
		element, err := typeCast(condition)
		if err != nil {
			return d, err
		}
		elements[condition.Field] = append(elements[condition.Field], element)
	}
	for key, values := range elements {
		d = append(d, bson.E{Key: key, Value: values})
	}
	return d, nil
}

// collections are the collection names from MongoDB but filtered so that not all collections
// are hammered every time we get events.
func (m *Database) collections(ctx context.Context, filter bson.D) ([]string, error) {
	value, ok := filter.Map()["meta.type"]
	// meta.type not set, return all collections.
	if !ok {
		return m.database.ListCollectionNames(ctx, bson.D{})
	}
	valueMap := value.(bson.D).Map()
	collection, ok := valueMap["$eq"]
	// No $eq. Apply collection filter to reduce the collection names
	// request.
	if !ok {
		// TODO: Collection filter
		return m.database.ListCollectionNames(ctx, bson.D{})
	}
	return []string{collection.(string)}, nil
}

// GetEvents gets all events information.
func (m *Database) GetEvents(ctx context.Context, request requests.MultipleEventsRequest) ([]drivers.EiffelEvent, error) {
	filter, err := buildFilter(request.Conditions)
	if err != nil {
		m.logger.Errorf("Database: %v", err)
		return nil, err
	}
	collections, err := m.collections(ctx, filter)
	if err != nil {
		m.logger.Errorf("Database: %v", err)
		return nil, err
	}

	m.logger.Debugf("fetching events from %d collections", len(collections))
	var allEvents []drivers.EiffelEvent
	for _, collection := range collections {
		var events []drivers.EiffelEvent
		cursor, err := m.database.Collection(collection).Find(ctx, filter,
			// Remove the _id field from the resulting document.
			options.Find().SetProjection(bson.M{"_id": 0}))
		if err != nil {
			continue
		}
		if err = cursor.All(ctx, &events); err != nil {
			m.logger.Info(err.Error())
			continue
		}
		allEvents = append(allEvents, events...)
	}
	return allEvents, nil
}

// SearchEvent searches for an event based on event ID.
func (m *Database) SearchEvent(ctx context.Context, id string) (drivers.EiffelEvent, error) {
	return drivers.EiffelEvent{}, errors.New("not yet implemented")
}

// UpstreamDownstreamSearch searches for events upstream and/or downstream of event by ID.
func (m *Database) UpstreamDownstreamSearch(ctx context.Context, id string) ([]drivers.EiffelEvent, error) {
	return nil, errors.New("not yet implemented")
}

// GetEventByID gets an event by ID in all collections.
func (m *Database) GetEventByID(ctx context.Context, id string) (drivers.EiffelEvent, error) {
	collections, err := m.collections(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	filter := bson.D{{Key: "meta.id", Value: id}}
	for _, collection := range collections {
		var event bson.M
		singleResult := m.database.Collection(collection).FindOne(ctx, filter,
			// Remove the _id field from the resulting document.
			options.FindOne().SetProjection(bson.M{"_id": 0}))
		err := singleResult.Decode(&event)
		if err != nil {
			continue
		} else {
			return drivers.EiffelEvent(event), nil
		}
	}
	return nil, fmt.Errorf("%q not found in any collection", id)
}

// Close the database connection.
func (m *Database) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
