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

// requests is a place for request structs to be used for both
// database and handlers.
package requests

import "github.com/eiffel-community/eiffel-goer/internal/query"

type MultipleEventsRequest struct {
	Shallow       bool  `schema:"shallow"` // TODO: Unused
	PageNo        int32 `schema:"pageNo"`
	PageSize      int32 `schema:"pageSize"`
	PageStartItem int32 `schema:"pageStartItem"`
	Lazy          bool  `schema:"lazy"`
	Readable      bool  `schema:"readable"` // TODO: Unused
	Conditions    []query.Condition
}

type SingleEventRequest struct {
	Shallow bool `schema:"shallow"` // TODO: Unused
}
