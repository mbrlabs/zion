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

package zion

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Context #TODO
type Context struct {
	zion *Zion

	Writer    http.ResponseWriter
	Request   *http.Request
	URLParams map[string]string
	User      User
}

// NewContext #TODO
func NewContext(z *Zion, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		zion:      z,
		Writer:    w,
		Request:   r,
		URLParams: make(map[string]string),
	}
}

// Render #TODO
func (ctx *Context) Render(name string, data interface{}) {
	ctx.zion.config.TemplateEngine.Render(name, data, ctx.Writer)
}

func (ctx *Context) Json(data interface{}) {
	data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("[ERROR] " + err.Error())
	}
	ctx.Writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(ctx.Writer, "%s", data)
}

func (ctx *Context) SendStatus(status int) {
	ctx.Writer.WriteHeader(status)
}

func (ctx *Context) Redirect(path string) {
	http.Redirect(ctx.Writer, ctx.Request, path, http.StatusTemporaryRedirect)
}
