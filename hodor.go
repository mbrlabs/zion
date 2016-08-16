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

func (this *Hodor) MountAfter(pattern string, middleware Middleware) {
	this.router.mountAfter(pattern, middleware)
}

func (this *Hodor) MountBefore(pattern string, middleware Middleware) {
	this.router.mountBefore(pattern, middleware)
}

func (this *Hodor) Get(pattern string, handler HandlerFunc) {
	this.router.addRoute(pattern, http.MethodGet, handler)
}

func (this *Hodor) Head(pattern string, handler HandlerFunc) {
	this.router.addRoute(pattern, http.MethodHead, handler)
}

func (this *Hodor) Post(pattern string, handler HandlerFunc) {
	this.router.addRoute(pattern, http.MethodPost, handler)
}

func (this *Hodor) Put(pattern string, handler HandlerFunc) {
	this.router.addRoute(pattern, http.MethodPut, handler)
}

func (this *Hodor) Delete(pattern string, handler HandlerFunc) {
	this.router.addRoute(pattern, http.MethodDelete, handler)
}

func (this *Hodor) Options(pattern string, handler HandlerFunc) {
	this.router.addRoute(pattern, http.MethodOptions, handler)
}

func (this *Hodor) configServer() {
	this.server.Addr = this.Host + ":" + strconv.Itoa(this.Port)
	this.server.ReadTimeout = this.ReadTimeout
	this.server.WriteTimeout = this.WriteTimeout
	this.server.Handler = this.router
}

func (this *Hodor) Start() {
	fmt.Println("Listening on http://localhost:3000")

	this.configServer()
	this.server.ListenAndServe()
}
