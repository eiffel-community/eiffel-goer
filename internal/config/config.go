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
package config

import (
	"flag"
	"os"
)

type Config interface {
	DBConnectionString() string
	APIPort() string
	LogLevel() string
	LogFilePath() string
}

type Cfg struct {
	connectionString string
	apiPort          string
	logLevel         string
	logFilePath      string
}

// Get parses input parameters to program and return a config with them set.
func Get() Config {
	conf := &Cfg{}

	flag.StringVar(&conf.connectionString, "connectionstring", os.Getenv("CONNECTION_STRING"), "Database connection string.")
	flag.StringVar(&conf.apiPort, "apiport", os.Getenv("API_PORT"), "API port.")
	flag.StringVar(&conf.logLevel, "loglevel", os.Getenv("LOGLEVEL"), "Log level (TRACE, DEBUG, INFO, WARNING, ERROR, FATAL, PANIC).")
	flag.StringVar(&conf.logFilePath, "logfilepath", os.Getenv("LOG_FILE_PATH"), "Path, including filename, for the log files to create.")

	flag.Parse()
	return conf
}

// DBConnectionString returns the connection string for a database.
func (c *Cfg) DBConnectionString() string {
	return c.connectionString
}

// APIPort returns the API port with a ":" prepended.
func (c *Cfg) APIPort() string {
	return ":" + c.apiPort
}

// LogLevel returns the log level. Default is INFO.
func (c *Cfg) LogLevel() string {
	if c.logLevel == "" {
		c.logLevel = "INFO"
	}
	return c.logLevel
}

// LogFilePath returns the path to where log files should be stored, including filename.
func (c *Cfg) LogFilePath() string {
	return c.logFilePath
}
