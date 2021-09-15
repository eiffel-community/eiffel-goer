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
package query

import (
	"net/url"
	"strings"

	"github.com/gorilla/schema"
	log "github.com/sirupsen/logrus"
)

// Query is a structure for parsing the query string in a way that
// is compatible with the API specification for an event repository.
type Query struct {
	URL    *url.URL
	Logger *log.Entry
}

// New creates a new Query parser for a URL.
func New(url *url.URL, logger *log.Entry) *Query {
	return &Query{
		url,
		logger,
	}
}

// Decode will decode a URL Query into a struct of 'schema' tagged values.
func (q *Query) Decode(dst interface{}) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	return decoder.Decode(dst, q.URL.Query())
}

// operators is a slice of operators that are allowed in params.
// %3C = <
// %3E = >
var operators = []string{"!=", "%3C=", "%3E=", "%3C", "%3E", "="}

// DecodeFilterParameters will decode the special search parameters of the
// ER specification. These parameters do not follow HTTP and need their own
// decoder. This decoder works very similar to url.URL.Query().
//
// From specification:
// `<resource>?key[.key ...]<FC>value[&key[.key ...]<FC>value ...]`
//
// <FC> is one of the Filter comparators described below. To traverse into nested structures and filter on their keys, namespacing with
// `.` (dot) is used.
//
// ```
// =  - equal to
// >  - greater than
// <  - less than
// >= - greater than or equal to
// <= - less than or equal to
// != - not equal to
// ```
func (q *Query) DecodeFilterParameters(ignoreKeys map[string]struct{}) (Params, error) {
	var err error
	params := Params{}
	rawQuery := q.URL.RawQuery
	for rawQuery != "" {
		key := rawQuery
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, rawQuery = key[:i], key[i+1:]
		} else {
			rawQuery = ""
		}
		if key == "" {
			continue
		}
		value := ""
		op := ""
		var ignore bool
		for _, operator := range operators {
			if i := strings.Index(key, operator); i >= 0 {
				op, err = url.QueryUnescape(operator)
				if err != nil {
					continue
				}
				key, value = key[:i], key[i+len(operator):]
				_, ok := ignoreKeys[key]
				if ok {
					ignore = true
				}
				break
			} else {
				continue
			}
		}
		if ignore {
			continue
		}
		key, err1 := url.QueryUnescape(key)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		value, err1 = url.QueryUnescape(value)
		if err1 != nil {
			if err == nil {
				err = err1
			}
			continue
		}
		params.Add(op, key, value)
	}
	return params, err
}

type Params map[string][][2]string

// Add appends value to key based on which operator was in query.
func (p Params) Add(operator string, key string, value string) {
	p[key] = append(p[key], [2]string{operator, value})
}
