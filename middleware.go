package hodor

type Middleware interface {
	Execute(ctx *Context) bool
	Name() string
}
