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
package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test that it is possible to get a Cfg from Get with values taken from environment variables.
func TestGet(t *testing.T) {
	port := "8080"
	connectionString := "connection string"
	logLevel := "DEBUG"
	logFilePath := "path/to/a/file"
	os.Setenv("CONNECTION_STRING", connectionString)
	os.Setenv("API_PORT", port)
	os.Setenv("LOGLEVEL", logLevel)
	os.Setenv("LOG_FILE_PATH", logFilePath)

	cfg, ok := Get().(*Cfg)
	assert.Truef(t, ok, "cfg returned from get is not a config interface")
	assert.Equal(t, connectionString, cfg.connectionString)
	assert.Equal(t, port, cfg.apiPort)
	assert.Equal(t, logLevel, cfg.logLevel)
	assert.Equal(t, logFilePath, cfg.logFilePath)
}

type getter func() string

// Test that the getters in the Cfg struct return the values from the struct.
func TestGetters(t *testing.T) {
	cfg := &Cfg{
		connectionString: "something://db/test",
		apiPort:          "8080",
		logLevel:         "TRACE",
		logFilePath:      "a/file/path.json",
	}
	emptyCfg := &Cfg{}
	tests := []struct {
		name     string
		cfg      *Cfg
		function getter
		value    string
	}{
		{name: "DBConnectionString", cfg: cfg, function: cfg.DBConnectionString, value: cfg.connectionString},
		{name: "APIPort", cfg: cfg, function: cfg.APIPort, value: ":" + cfg.apiPort},
		{name: "LogLevel", cfg: cfg, function: cfg.LogLevel, value: cfg.logLevel},
		{name: "LogLevelDefault", cfg: emptyCfg, function: emptyCfg.LogLevel, value: "INFO"},
		{name: "LogFilePath", cfg: cfg, function: cfg.LogFilePath, value: cfg.logFilePath},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.value, testCase.function())
		})
	}
}
