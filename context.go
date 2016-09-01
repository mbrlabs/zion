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
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
)

// extra keys
const (
	ExtraUser = "extra_user"
)

// Context
//------------------------------------------------------------------------------------

// Context #TODO
type Context struct {
	zion *Zion

	writer  http.ResponseWriter
	request *http.Request

	urlParams map[string]string
	extras    map[string]interface{}
}

// NewContext #TODO
func NewContext(z *Zion, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		zion:      z,
		writer:    w,
		request:   r,
		urlParams: make(map[string]string),
		extras:    make(map[string]interface{}),
	}
}

// Render #TODO
func (ctx *Context) Render(name string, data interface{}) {
	ctx.zion.config.TemplateEngine.Render(name, data, ctx.writer)
}

func (ctx *Context) Html(html string) {
	fmt.Fprintf(ctx.writer, "%s", html)
}

func (ctx *Context) Json(data interface{}) {
	data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("[ERROR] " + err.Error())
	}
	ctx.writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(ctx.writer, "%s", data)
}

// TODO this sends the file EVERY request. User Last-Modified instead and send only if file changed.
func (ctx *Context) File(path string) {
	fmt.Println("Context: serving file " + path)

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		ctx.SendStatus(404)
		return
	}

	stats, _ := file.Stat()
	ctx.writer.Header().Set("Content-Disposition", "attachment; filename="+file.Name())
	ctx.writer.Header().Set("Content-Type", mime.TypeByExtension(path))
	ctx.writer.Header().Set("Content-Length", strconv.FormatInt(stats.Size(), 10))

	io.Copy(ctx.writer, file)
}

func (ctx *Context) SendStatus(status int) {
	ctx.writer.WriteHeader(status)
}

func (ctx *Context) Redirect(path string) {
	http.Redirect(ctx.writer, ctx.request, path, http.StatusTemporaryRedirect)
}

func (ctx *Context) Method() string {
	return ctx.request.Method
}

func (ctx *Context) Path() string {
	return ctx.request.URL.Path
}

func (ctx *Context) URLParams() map[string]string {
	return ctx.urlParams
}

func (ctx *Context) URLParam(name string) string {
	return ctx.urlParams[name]
}

func (ctx *Context) FormValue(name string) string {
	return ctx.request.FormValue(name)
}

func (ctx *Context) Cookie(name string) (*http.Cookie, error) {
	return ctx.request.Cookie(name)
}

func (ctx *Context) SetCookie(c *http.Cookie) {
	http.SetCookie(ctx.writer, c)
}

func (ctx *Context) Extra(key string) interface{} {
	return ctx.extras[key]
}

func (ctx *Context) AddExtra(key string, extra interface{}) {
	ctx.extras[key] = extra
}
