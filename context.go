package hodor

import (
	"net/http"
)

type Context struct {
	Writer  http.ResponseWriter
	Request *http.Request

	UrlParams map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:    w,
		Request:   r,
		UrlParams: make(map[string]string),
	}
}
