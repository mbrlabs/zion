package hodor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Hodor #TODO
type Hodor struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	server http.Server
	router *Router

	templateEngine HTMLTemplateEngine
}

// NewHodor returns a new Hodor instance
func (h *Hodor) configServer() {
	h.server.Addr = h.Host + ":" + strconv.Itoa(h.Port)
	h.server.ReadTimeout = h.ReadTimeout
	h.server.WriteTimeout = h.WriteTimeout
	h.server.Handler = h.router
}

// NewHodor #TODO
func NewHodor() *Hodor {
	app := &Hodor{
		Port:           3000,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		templateEngine: NewDefaultTemplateEngine(),
	}
	app.router = NewRouter(app)

	return app
}

// MountAfter #TODO
func (h *Hodor) MountAfter(pattern string, middleware Middleware) {
	h.router.mountAfter(pattern, middleware)
}

// MountBefore #TODO
func (h *Hodor) MountBefore(pattern string, middleware Middleware) {
	h.router.mountBefore(pattern, middleware)
}

// Get #TODO
func (h *Hodor) Get(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodGet, handler)
}

// Head #TODO
func (h *Hodor) Head(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodHead, handler)
}

// Post #TODO
func (h *Hodor) Post(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPost, handler)
}

// Put #TODO
func (h *Hodor) Put(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodPut, handler)
}

// Delete #TODO
func (h *Hodor) Delete(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodDelete, handler)
}

// Options #TODO
func (h *Hodor) Options(pattern string, handler HandlerFunc) {
	h.router.addRoute(pattern, http.MethodOptions, handler)
}

// ServeStaticFiles #TODO
func (h *Hodor) ServeStaticFiles(urlPath string, fsPath string) {
	staticPrefix := strings.Trim(urlPath, "/")
	if strings.HasSuffix(staticPrefix, "*") {
		staticPrefix = strings.TrimRight(staticPrefix, "*")
		staticPrefix = "/" + strings.Trim(staticPrefix, "/") + "/"

		// Server files
		fileServer := http.StripPrefix(staticPrefix, http.FileServer(http.Dir(fsPath)))
		h.Get(urlPath, func(ctx *Context) {
			fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	} else {
		panic("Static files must be mapped with a wildcard (*) in the pattern url")
	}
}

func (h *Hodor) Templates(path string) {
	h.templateEngine.CompileTemplates(path)
}

// Start #TODO
func (h *Hodor) Start() {
	h.configServer()
	fmt.Println("Listening on http://localhost:3000")

	h.server.ListenAndServe()
}
