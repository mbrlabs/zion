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
	"strconv"
	"strings"
)

// Zion #TODO
type Zion struct {
	server http.Server
	router *router
	config *Config
}

// New #TODO
func New(config *Config) *Zion {
	app := &Zion{config: config}
	app.router = newRouter(app)

	return app
}

// MountAfter #TODO
func (z *Zion) MountAfter(pattern string, middleware Middleware) {
	z.router.mountAfter(pattern, middleware)
}

// MountBefore #TODO
func (z *Zion) MountBefore(pattern string, middleware Middleware) {
	z.router.mountBefore(pattern, middleware)
}

// Get #TODO
func (z *Zion) Get(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodGet, handler)
}

// Head #TODO
func (z *Zion) Head(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodHead, handler)
}

// Post #TODO
func (z *Zion) Post(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodPost, handler)
}

// Put #TODO
func (z *Zion) Put(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodPut, handler)
}

// Delete #TODO
func (z *Zion) Delete(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodDelete, handler)
}

// Options #TODO
func (z *Zion) Options(pattern string, handler HandlerFunc) {
	z.router.addRoute(pattern, http.MethodOptions, handler)
}

// ServeStaticFiles #TODO
func (z *Zion) ServeStaticFiles(urlPath string, fsPath string) {
	staticPrefix := strings.Trim(urlPath, "/")
	if strings.HasSuffix(staticPrefix, "*") {
		staticPrefix = strings.TrimRight(staticPrefix, "*")
		staticPrefix = "/" + strings.Trim(staticPrefix, "/") + "/"

		// Server files
		fileServer := http.StripPrefix(staticPrefix, http.FileServer(http.Dir(fsPath)))
		z.Get(urlPath, func(ctx *Context) {
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	} else {
		panic("Static files must be mapped with a wildcard (*) in the pattern url")
	}
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
