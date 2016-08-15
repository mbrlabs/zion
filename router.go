package hodor

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(ctx *Context)

// ============================================================================
// 								struct Route
// ============================================================================
type Route struct {
	path    string
	Method  string
	Handler HandlerFunc
}

func NewRoute(pattern string, method string, handler HandlerFunc) Route {
	route := Route{Method: method, Handler: handler}
	route.SetPath(pattern)
	return route
}

func (r *Route) SetPath(pattern string) {
	r.path = strings.Trim(pattern, "/")
}

func (r *Route) GetPath() string {
	return r.path
}

// ============================================================================
// 								struct Router
// ============================================================================
type Router struct {
	routes []Route
	after  []Middleware
	before []Middleware
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
	r.routes = append(r.routes, NewRoute(pattern, method, handler))
}

// TODO replace this with the suffix tree implementation
func (r *Router) findRoute(pattern string) *Route {
	pattern = strings.Trim(pattern, "/")
	for _, route := range r.routes {
		if route.GetPath() == pattern {
			return &route
		}
	}

	return nil
}

func (r Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := &Context{Writer: resp, Request: req}

	route := r.findRoute(req.URL.Path)
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
		route.Handler(ctx)
		// after middleware
		for _, mw := range r.after {
			if !mw.Execute(ctx) {
				return
			}
		}
	}
}
