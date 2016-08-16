package hodor

import (
    "testing"
)

func TestBasics(t *testing.T) {
    tree := newRouteTree()

    tree.insert(newRoute("/", "GET", func(ctx *Context) {}))
    tree.insert(newRoute("/user", "GET", func(ctx *Context) {}))
    tree.insert(newRoute("/user/:param", "GET", func(ctx *Context) {}))
    tree.insert(newRoute("/home/:param", "GET", func(ctx *Context) {}))
    tree.insert(newRoute("/:param/hannah", "GET", func(ctx *Context) {}))
    tree.insert(newRoute("/:param/hannah/:page", "GET", func(ctx *Context) {}))
}