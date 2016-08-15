package hodor

import (
	"net/http"
	"fmt"
)

type Router struct {
	routes 	[]Route
	after	[]Middleware
	before	[]Middleware
}

type Route struct {
	path 	string
	method 	string
	handler func(ctx *Context)
}

func (r *Router) mountAfter(pattern string, middleware Middleware) {
	r.after = append(r.after, middleware)
	fmt.Printf("Added middleware %d\n", len(r.after))
}

func (r *Router) mountBefore(pattern string, middleware Middleware) {
	r.before = append(r.before, middleware)
	fmt.Printf("Added middleware %d\n", len(r.before))
}

func (r *Router) addRoute(pattern string, method string, handler func(ctx *Context)) {
	r.routes = append(r.routes, Route{path: pattern, method: method, handler: handler})
}

func (r *Router) findRoute(pattern string) *Route {
	for _, route := range r.routes {
		if route.path == pattern {
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
			mw.Execute(ctx)
		}

		// actual route
		route.handler(ctx)

		// after middleware
		for _, mw := range r.after {
			mw.Execute(ctx)
		}
	}
}