package hodor

import (
	"net/http"
	"fmt"
)

type Router struct {
	routes 	map[string]func(http.ResponseWriter, *http.Request)
	after	[]Middleware
	before	[]Middleware
}

type Route struct {
	path 	string
	method 	string
	handler func(http.ResponseWriter, *http.Request)
}

func NewRouter() *Router {
	return &Router {
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (r *Router) mountAfter(pattern string, middleware Middleware) {
	r.after = append(r.after, middleware)
	fmt.Printf("Added middleware %d\n", len(r.after))
}

func (r *Router) mountBefore(pattern string, middleware Middleware) {
	r.before = append(r.before, middleware)
	fmt.Printf("Added middleware %d\n", len(r.before))
}

func (r *Router) addRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.routes[pattern] = handler
}

func (r Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {	
	ctx := &Context{Writer: resp, Request: req}

	handler := r.routes[req.URL.Path]
	if handler == nil {
		// we didn't find a handler -> send a 404
		http.NotFound(resp, req)
	} else {

		// before middleware
		for _, mw := range r.before {
			mw.Execute(ctx)
		}

		// actual route
		handler(resp, req)

		// after middleware
		for _, mw := range r.after {
			mw.Execute(ctx)
		}
	}
}