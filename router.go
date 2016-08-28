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
	"runtime/debug"
	"strings"
)

const (
	pageNotFoundPage = "<!DOCTYPE html><html><head><title>Page Not Found</title></head><body>" +
		"<h2>Page not found</h2>" +
		"</body></html>"
	serverErrorPageProduction = "<!DOCTYPE html><html><head><title>Server Error</title></head><body>" +
		"<h2>Internal Server Error</h2>" +
		"</body></html>"
	serverErrorPageDevelopment = "<!DOCTYPE html><html><head><title>Server Error</title></head><body>" +
		"<h2>Internal Server Error</h2></br>%s" +
		"</body></html>"
)

// HandlerFunc #TODO
type HandlerFunc func(ctx *Context)

// ============================================================================
// 								struct route
// ============================================================================
type route struct {
	pattern    string
	method     string
	handler    HandlerFunc
	middleware []Middleware
}

func newRoute(pattern string, method string, handler HandlerFunc) *route {
	r := &route{method: method, handler: handler}
	r.setPattern(pattern)
	return r
}

func (r *route) setPattern(pattern string) {
	r.pattern = strings.Trim(pattern, "/")
}

func (r *route) getPattern() string {
	return r.pattern
}

// ============================================================================
// 								struct Router
// ============================================================================

// Router #TODO
type router struct {
	tree   routeTree
	after  []Middleware
	before []Middleware

	zion *Zion
}

// NewRouter #TODO
func newRouter(zion *Zion) *router {
	return &router{zion: zion, tree: newRouteTree()}
}

func (r *router) mountAfter(pattern string, middleware Middleware) {
	r.after = append(r.after, middleware)
	fmt.Printf("mountAfter: %s\n", middleware.Name())
}

func (r *router) mountBefore(pattern string, middleware Middleware) {
	r.before = append(r.before, middleware)
	fmt.Printf("mountBefore: %s\n", middleware.Name())
}

func (r *router) addRoute(pattern string, method string, handler HandlerFunc) {
	r.tree.insertRoute(newRoute(pattern, method, handler))
}

func (r *router) recover(w http.ResponseWriter, req *http.Request) {
	if obj := recover(); obj != nil {
		stacktrace := string(debug.Stack()[:])
		fmt.Printf("\n\n[CAPTURED PANIC] ====> \n\n%s\n\n[STACTRACE END] <====\n", stacktrace)
		if r.zion.config.DevelopmentMode {
			fmt.Fprintf(w, serverErrorPageDevelopment, escapeHTML(stacktrace))
		} else {
			fmt.Fprintf(w, serverErrorPageProduction)
		}
	}
}

func (r *router) serve(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext(r.zion, resp, req)
	route := r.tree.get(ctx)

	if route == nil {
		// we didn't find a handler -> send a 404 or redirect to error page
		if len(r.zion.config.PageNotFoundRedirect) > 0 {
			ctx.Redirect(r.zion.config.PageNotFoundRedirect)
		} else {
			http.NotFound(resp, req)
		}
	} else {
		// before middleware
		for _, mw := range r.before {
			if !mw.Execute(ctx) {
				return
			}
		}
		// actual route
		route.handler(ctx)
		// after middleware
		for _, mw := range r.after {
			if !mw.Execute(ctx) {
				return
			}
		}
	}
}

func (r *router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer r.recover(resp, req)
	r.serve(resp, req)
}
