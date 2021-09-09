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
package main

import (
	"context"
	"os"

	"github.com/eiffel-community/eiffel-goer/internal/config"
	"github.com/eiffel-community/eiffel-goer/internal/logger"
	"github.com/eiffel-community/eiffel-goer/pkg/application"
	log "github.com/sirupsen/logrus"
)

// GitSummary contains "git describe" output and is automatically
// populated via linker options when building with govvv.
var GitSummary = "(unknown)"

// Start up the Goer application.
func main() {
	cfg := config.Get()
	ctx := context.Background()
	if err := logger.Setup(cfg); err != nil {
		log.Fatal(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	log := log.WithFields(log.Fields{
		"hostname":    hostname,
		"application": "eiffel-goer",
		"version":     GitSummary,
	})
	app, err := application.Get(ctx, cfg, log)
	if err != nil {
		log.Panic(err)
	}

	app.LoadV1Alpha1Routes()

	log.Debug("Starting up.")
	err = app.Start(ctx)
	if err != nil {
		log.Panic(err)
	}
}
