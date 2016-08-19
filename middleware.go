package hodor

// Middleware #TODO
type Middleware interface {
	Execute(ctx *Context) bool
	Name() string
}
