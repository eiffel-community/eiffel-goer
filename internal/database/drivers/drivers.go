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
package drivers

import (
	"context"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/eiffel-community/eiffel-goer/internal/schema"
)

type DatabaseDriver interface {
	Get(*url.URL, *log.Entry) (DatabaseDriver, error)
	SupportsScheme(string) bool
	Connect(context.Context) error
	GetEvents(context.Context) ([]schema.EiffelEvent, error)
	SearchEvent(context.Context, string) (schema.EiffelEvent, error)
	UpstreamDownstreamSearch(context.Context, string) ([]schema.EiffelEvent, error)
	GetEventByID(context.Context, string) (schema.EiffelEvent, error)
	Close(context.Context) error
}