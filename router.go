package hodor

import (
	"fmt"
	"net/http"
	"strings"
)

// HandlerFunc #TODO
type HandlerFunc func(ctx *Context)

// ============================================================================
// 								struct route
// ============================================================================
type route struct {
	pattern string
	method  string
	handler HandlerFunc
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
	r.tree.insert(newRoute(pattern, method, handler))
}

func (r *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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
