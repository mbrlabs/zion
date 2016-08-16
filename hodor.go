package hodor

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type Hodor struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	server http.Server
	router *Router
}

// Returns a new Hodor instance
//
func NewHodor() *Hodor {
	app := &Hodor{
		Port:         3000,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		router:       &Router{},
	}

	return app
}

func (h *Hodor) MountAfter(pattern string, middleware Middleware) {
	h.router.mountAfter(pattern, middleware)
}

func (h *Hodor) MountBefore(pattern string, middleware Middleware) {
	h.router.mountBefore(pattern, middleware)
}

func (h *Hodor) Get(pattern string, handler func(ctx *Context)) {
	h.router.addRoute(pattern, "GET", handler)
}

func (h *Hodor) configServer() {
	h.server.Addr = h.Host + ":" + strconv.Itoa(h.Port)
	h.server.ReadTimeout = h.ReadTimeout
	h.server.WriteTimeout = h.WriteTimeout
}

func (h *Hodor) Start() {
	fmt.Println("Listening on http://localhost:3000")

	h.configServer()
	h.server.ListenAndServe()
}
