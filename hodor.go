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

package hodor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Hodor #TODO
type Hodor struct {
	server http.Server
	router *router
	config *Config
}

// NewHodor #TODO
func NewHodor(config *Config) *Hodor {
	app := &Hodor{config: config}
	app.router = newRouter(app)

	return app
}

// MountAfter #TODO
func (h *Hodor) MountAfter(pattern string, middleware Middleware) {
	h.router.mountAfter(pattern, middleware)
}

// MountBefore #TODO
func (h *Hodor) MountBefore(pattern string, middleware Middleware) {
	h.router.mountBefore(pattern, middleware)
}

// Get #TODO
func (h *Hodor) Get(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodGet, handler)
}

// Head #TODO
func (h *Hodor) Head(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodHead, handler)
}

// Post #TODO
func (h *Hodor) Post(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPost, handler)
}

// Put #TODO
func (h *Hodor) Put(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPut, handler)
}

// Delete #TODO
func (h *Hodor) Delete(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodDelete, handler)
}

// Options #TODO
func (h *Hodor) Options(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodOptions, handler)
}

// ServeStaticFiles #TODO
func (h *Hodor) ServeStaticFiles(urlPath string, fsPath string) {
	staticPrefix := strings.Trim(urlPath, "/")
	if strings.HasSuffix(staticPrefix, "*") {
		staticPrefix = strings.TrimRight(staticPrefix, "*")
		staticPrefix = "/" + strings.Trim(staticPrefix, "/") + "/"

		// Server files
		fileServer := http.StripPrefix(staticPrefix, http.FileServer(http.Dir(fsPath)))
		h.Get(urlPath, func(ctx *Context) {
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	} else {
		panic("Static files must be mapped with a wildcard (*) in the pattern url")
	}
}

func (h *Hodor) checkConfig() {
	// TODO check server timeouts

	// TODO check host + port config

	// TODO check template settings

	// TODO print development mode warnings
}

func (h *Hodor) prepare() {
	// prepare server
	h.server.Addr = h.config.Host + ":" + strconv.Itoa(h.config.Port)
	h.server.ReadTimeout = h.config.ReadTimeout
	h.server.WriteTimeout = h.config.WriteTimeout
	h.server.Handler = h.router

	// setup static files
	h.ServeStaticFiles(h.config.StaticFileURLPattern, h.config.StaticFilePath)

	// compile templates
	h.config.TemplateEngine.CompileTemplates(h.config.TemplatePath)

	// configure development mode stuff
	if h.config.DevelopmentMode {
		h.config.TemplateEngine.EnableRecompiling(true)
	}
}

// Start #TODO
func (h *Hodor) Start() {
	h.checkConfig()
	h.prepare()
	fmt.Println("Listening on http://" + h.config.Host + ":" + strconv.Itoa(h.config.Port))

	h.server.ListenAndServe()
}
