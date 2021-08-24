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

// This package is only for generating the mocks that are used in test.
package test

//go:generate mockgen -destination mock_database/mock_database.go github.com/eiffel-community/eiffel-goer/internal/database Database
//go:generate mockgen -destination mock_config/mock_config.go github.com/eiffel-community/eiffel-goer/internal/config Config
//go:generate mockgen -destination mock_server/mock_server.go github.com/eiffel-community/eiffel-goer/pkg/server Server
