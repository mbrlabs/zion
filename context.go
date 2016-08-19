package hodor

import (
	"net/http"
)

// Context #TODO
type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	URLParams map[string]string
}

// NewContext #TODO
func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:    w,
		Request:   r,
		URLParams: make(map[string]string),
	}
}
