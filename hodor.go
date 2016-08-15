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

func (h *Hodor) MountAfter(pattern string, middleware Middleware) {
	h.router.mountAfter(pattern, middleware)
}

func (h *Hodor) MountBefore(pattern string, middleware Middleware) {
	h.router.mountBefore(pattern, middleware)
}

func (h *Hodor) Get(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	h.router.addRoute(pattern, handler)
}

func (h *Hodor) Start() {
	fmt.Println("Listening on http://localhost:3000")

	h.server.Handler = *h.router
	h.server.ListenAndServe()
}