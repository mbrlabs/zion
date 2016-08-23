package hodor

import (
	"net/http"
)

// Context #TODO
type Context struct {
	hodor *Hodor

	Writer    http.ResponseWriter
	Request   *http.Request
	URLParams map[string]string
	User      *User
}

// NewContext #TODO
func NewContext(h *Hodor, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		hodor:     h,
		Writer:    w,
		Request:   r,
		URLParams: make(map[string]string),
	}
}

// Render #TODO
func (ctx *Context) Render(name string, data interface{}) {
	ctx.hodor.templateEngine.Render(name, data, ctx.Writer)
}
