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
	"runtime/debug"
	"strings"
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
type Router struct {
	tree   routeTree
	after  []Middleware
	before []Middleware

	hodor *Hodor
}

// NewRouter #TODO
func NewRouter(hodor *Hodor) *Router {
	return &Router{hodor: hodor, tree: newRouteTree()}
}

func (r *Router) mountAfter(pattern string, middleware Middleware) {
	r.after = append(r.after, middleware)
	fmt.Printf("mountAfter: %s\n", middleware.Name())
}

func (r *Router) mountBefore(pattern string, middleware Middleware) {
	r.before = append(r.before, middleware)
	fmt.Printf("mountBefore: %s\n", middleware.Name())
}

func (r *Router) addRoute(pattern string, method string, handler HandlerFunc) {
	r.tree.insertRoute(newRoute(pattern, method, handler))
}

func (r *Router) recover(resp http.ResponseWriter, req *http.Request) {
	if obj := recover(); obj != nil {
		fmt.Printf("\n\n[CAPTURED PANIC] ====> \n\n%s\n\n[STACTRACE END] <====\n", string(debug.Stack()[:]))
	}

	// TODO send depending on dev mode a detailed message.
	// also make it possible to define a custom page in case of a panic recovery.

	// send internal server error
	resp.WriteHeader(http.StatusInternalServerError)
}

func (r *Router) serve(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext(r.hodor, resp, req)
	route := r.tree.get(ctx)

	if route == nil {
		// we didn't find a handler -> send a 404
		http.NotFound(resp, req)
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

func (r *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer r.recover(resp, req)
	r.serve(resp, req)
}
