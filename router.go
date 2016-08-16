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
type route struct {
	pattern string
	method  string
	handler HandlerFunc
}

func newRoute(pattern string, method string, handler HandlerFunc) route {
	r := route{method: method, handler: handler}
	r.setPattern(pattern)
	return r
}

func (this *route) setPattern(pattern string) {
	this.pattern = strings.Trim(pattern, "/")
}

func (this *route) getPattern() string {
	return this.pattern
}

// ============================================================================
// 								struct Router
// ============================================================================
type Router struct {
	tree   routeTree
	after  []Middleware
	before []Middleware
}

func NewRouter() *Router {
	return &Router{tree: newRouteTree()}
}

func (this *Router) mountAfter(pattern string, middleware Middleware) {
	this.after = append(this.after, middleware)
	fmt.Printf("mountAfter: %s\n", middleware.Name())
}

func (this *Router) mountBefore(pattern string, middleware Middleware) {
	this.before = append(this.before, middleware)
	fmt.Printf("mountBefore: %s\n", middleware.Name())
}

func (this *Router) addRoute(pattern string, method string, handler HandlerFunc) {
	route := newRoute(pattern, method, handler)
	this.tree.insert(&route)
}

func (this Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := NewContext(resp, req)
	route := this.tree.get(ctx)

	if route == nil {
		// we didn't find a handler -> send a 404
		http.NotFound(resp, req)
	} else {
		// before middleware
		for _, mw := range this.before {
			if !mw.Execute(ctx) {
				return
			}
		}
		// actual route
		route.handler(ctx)
		// after middleware
		for _, mw := range this.after {
			if !mw.Execute(ctx) {
				return
			}
		}
	}
}
