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
		router:       NewRouter(),
	}

	return app
}

func (h *Hodor) MountAfter(pattern string, middleware Middleware) {
	h.router.mountAfter(pattern, middleware)
}

func (h *Hodor) MountBefore(pattern string, middleware Middleware) {
	h.router.mountBefore(pattern, middleware)
}

func (h *Hodor) Get(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodGet, handler)
}

func (h *Hodor) Head(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodHead, handler)
}

func (h *Hodor) Post(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPost, handler)
}

func (h *Hodor) Put(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPut, handler)
}

func (h *Hodor) Patch(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPatch, handler)
}

func (h *Hodor) Delete(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodDelete, handler)
}

func (h *Hodor) Connect(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodConnect, handler)
}

func (h *Hodor) Options(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodOptions, handler)
}

func (h *Hodor) Trace(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodTrace, handler)
}

func (h *Hodor) configServer() {
	h.server.Addr = h.Host + ":" + strconv.Itoa(h.Port)
	h.server.ReadTimeout = h.ReadTimeout
	h.server.WriteTimeout = h.WriteTimeout
	h.server.Handler = h.router
}

func (h *Hodor) Start() {
	fmt.Println("Listening on http://localhost:3000")

	h.configServer()
	h.server.ListenAndServe()
}
