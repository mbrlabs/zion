package hodor

import (
	"testing"
)

func TestBasics(t *testing.T) {
	tree := newRouteTree()

	tree.insertRoute(newRoute("/", "GET", func(ctx *Context) {}))
	tree.insertRoute(newRoute("/user", "GET", func(ctx *Context) {}))
	tree.insertRoute(newRoute("/user/:param", "GET", func(ctx *Context) {}))
	tree.insertRoute(newRoute("/home/:param", "GET", func(ctx *Context) {}))
	tree.insertRoute(newRoute("/:param/hannah", "GET", func(ctx *Context) {}))
	tree.insertRoute(newRoute("/:param/hannah/:page", "GET", func(ctx *Context) {}))
}
