// Copyright (c) 2016. See AUTHORS file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	User      User
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
