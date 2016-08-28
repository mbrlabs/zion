// Copyright (c) 2016. See AUTHORS file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zion

import (
	"time"
)

type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	TemplatePath   string
	TemplateEngine HTMLTemplateEngine

	StaticFilePath       string
	StaticFileURLPattern string

	PageNotFoundRedirect string
	ServerErrorRedirect  string

	DevelopmentMode bool
}

func NewConfig() *Config {
	return &Config{
		Host:                 "localhost",
		Port:                 3000,
		ReadTimeout:          10 * time.Second,
		WriteTimeout:         10 * time.Second,
		TemplatePath:         "views/",
		TemplateEngine:       NewDefaultTemplateEngine(),
		StaticFilePath:       "static/",
		StaticFileURLPattern: "/static/*",
		DevelopmentMode:      true,
	}
}
