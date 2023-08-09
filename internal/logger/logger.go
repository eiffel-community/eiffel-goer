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
package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"

	"github.com/eiffel-community/eiffel-goer/internal/config"
)

// Setup sets up logging to file with a JSON format and to stdout in text format.
func Setup(cfg config.Config) error {
	logLevel, err := logrus.ParseLevel(cfg.LogLevel())
	if err != nil {
		return err
	}
	if filePath := cfg.LogFilePath(); filePath != "" {
		// TODO: Make these parameters configurable.
		rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
			Filename:   filePath,
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     0, // days
			Level:      logrus.DebugLevel,
			Formatter:  &logrus.JSONFormatter{},
		})
		if err != nil {
			return err
		}
		logrus.AddHook(rotateFileHook)
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	return nil
}
