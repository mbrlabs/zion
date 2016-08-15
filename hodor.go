package hodor

import (
	"net/http"
	"fmt"
)

type Hodor struct {
	server 				http.Server
	router 				*Router
}

// Returns a new Hodor instance
//
func NewHodor() *Hodor {
	return &Hodor {
		server: http.Server {Addr: "127.0.0.1:3000"},
		router: NewRouter(),
	}
}

// Adds a middleware
func (h *Hodor) Use(handler func(http.ResponseWriter, *http.Request)) {
	h.router.use(handler)
}

func (h *Hodor) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.router.addRoute(pattern, handler)
}

func (h *Hodor) Start() {
	fmt.Println("Starting server...")

	h.server.Handler = *h.router
	h.server.ListenAndServe()
}