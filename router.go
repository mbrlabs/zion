package hodor

import (
	"net/http"
	"fmt"
)

type Router struct {
	routes 	map[string]func(http.ResponseWriter, *http.Request)
	middleware	[]func(http.ResponseWriter, *http.Request)
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

func (r *Router) use(handler func(http.ResponseWriter, *http.Request)) {
	r.middleware = append(r.middleware, handler)
	fmt.Printf("Added middleware %d\n", len(r.middleware))
}

func (r *Router) addRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.routes[pattern] = handler
}

func (r Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {	
	handler := r.routes[req.URL.Path]
	if handler == nil {
		// we didn't find a handler -> send a 404
		http.NotFound(resp, req)
	} else {
		// we found a matching handler -> go through middleware and then call the handler
		for _, mw := range r.middleware {
			mw(resp, req)
		}
		handler(resp, req)
	}
}