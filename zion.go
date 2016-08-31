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
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	HTTPStandard = "net/http"
	HTTPFast     = "valyala/fasthttp"
)

// Config
//------------------------------------------------------------------------------------

// Config is used to configure zion
type Config struct {
	HTTPImplementation string
	Host               string
	Port               int
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration

	TemplatePath   string
	TemplateEngine HTMLTemplateEngine

	StaticFilePath       string
	StaticFileURLPattern string

	PageNotFoundRedirect string
	ServerErrorRedirect  string

	DevelopmentMode bool
}

// NewConfig returns a new config with ready to use default values
func NewConfig() *Config {
	return &Config{
		Host:                 "localhost",
		Port:                 3000,
		ReadTimeout:          10 * time.Second,
		WriteTimeout:         10 * time.Second,
		TemplatePath:         "views/",
		TemplateEngine:       NewDefaultTemplateEngine(),
		StaticFilePath:       "static/",
		StaticFileURLPattern: "/static/",
		DevelopmentMode:      true,
		HTTPImplementation:   HTTPStandard,
	}
}

// Zion
//------------------------------------------------------------------------------------

// Zion is holds everythin together
type Zion struct {
	server http.Server
	router *router
	config *Config
}

// New returns a new zion instance
func New(config *Config) *Zion {
	app := &Zion{config: config}
	app.router = newRouter(app)

	return app
}

// MountAfter mounts middleware after the handler functions
func (z *Zion) MountAfter(pattern string, middleware Middleware) {
	z.router.mountAfter(pattern, middleware)
}

// MountBefore mounts middleware before the handler functions
func (z *Zion) MountBefore(pattern string, middleware Middleware) {
	z.router.mountBefore(pattern, middleware)
}

// Get registers the handler with the given url pattern for GET requests
func (z *Zion) Get(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodGet, handler)
}

// Head registers the handler with the given url pattern for HEAD requests
func (z *Zion) Head(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodHead, handler)
}

// Post registers the handler with the given url pattern for POST requests
func (z *Zion) Post(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodPost, handler)
}

// Put registers the handler with the given url pattern for PUT requests
func (z *Zion) Put(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodPut, handler)
}

// Delete registers the handler with the given url pattern for DELETE requests
func (z *Zion) Delete(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodDelete, handler)
}

// Options registers the handler with the given url pattern for OPTIONS requests
func (z *Zion) Options(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodOptions, handler)
}

// ServeStaticFiles serves static files
func (z *Zion) ServeStaticFiles(urPrefix string, fsPath string) {
	pattern := strings.Trim(urPrefix, "/") + "/*file"
	z.Get(pattern, func(ctx Context) {
		ctx.File(path.Join(fsPath, ctx.URLParams()["file"]))
	})
}

func (z *Zion) checkConfig() {
	// TODO check server timeouts

	// TODO check host + port config

	// TODO check template settings

	// TODO print development mode warnings
}

func (z *Zion) prepare() {
	// prepare server
	z.server.Addr = z.config.Host + ":" + strconv.Itoa(z.config.Port)
	z.server.ReadTimeout = z.config.ReadTimeout
	z.server.WriteTimeout = z.config.WriteTimeout
	z.server.Handler = z.router

	// setup static files
	z.ServeStaticFiles(z.config.StaticFileURLPattern, z.config.StaticFilePath)

	// compile templates
	z.config.TemplateEngine.CompileTemplates(z.config.TemplatePath)

	// configure development mode stuff
	if z.config.DevelopmentMode {
		z.config.TemplateEngine.EnableRecompiling(true)
	}
}

// Start #TODO
func (z *Zion) Start() {
	z.checkConfig()
	z.prepare()
	fmt.Println("Listening on http://" + z.config.Host + ":" + strconv.Itoa(z.config.Port))

	z.server.ListenAndServe()
}
