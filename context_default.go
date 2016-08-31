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

// Context
//------------------------------------------------------------------------------------

// contextDefault #TODO
type contextDefault struct {
	zion *Zion

	writer  http.ResponseWriter
	request *http.Request

	urlParams map[string]string

	extras map[string]interface{}
}

// NewContext #TODO
func newDefaultContext(z *Zion, w http.ResponseWriter, r *http.Request) Context {
	return &contextDefault{
		zion:      z,
		writer:    w,
		request:   r,
		urlParams: make(map[string]string),
		extras:    make(map[string]interface{}),
	}
}

// Render #TODO
func (ctx *contextDefault) Render(name string, data interface{}) {
	ctx.zion.config.TemplateEngine.Render(name, data, ctx.writer)
}

func (ctx *contextDefault) Html(html string) {
	fmt.Fprintf(ctx.writer, "%s", html)
}

func (ctx *contextDefault) Json(data interface{}) {
	data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("[ERROR] " + err.Error())
	}
	ctx.writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(ctx.writer, "%s", data)
}

// TODO this sends the file EVERY request. User Last-Modified instead and send only if file changed.
func (ctx *contextDefault) File(path string) {
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

func (ctx *contextDefault) SendStatus(status int) {
	ctx.writer.WriteHeader(status)
}

func (ctx *contextDefault) Redirect(path string) {
	http.Redirect(ctx.writer, ctx.request, path, http.StatusTemporaryRedirect)
}

func (ctx *contextDefault) Method() string {
	return ctx.request.Method
}

func (ctx *contextDefault) Path() string {
	return ctx.request.URL.Path
}

func (ctx *contextDefault) URLParams() map[string]string {
	return ctx.urlParams
}

func (ctx *contextDefault) URLParam(name string) string {
	return ctx.urlParams[name]
}

func (ctx *contextDefault) FormValue(name string) string {
	return ctx.request.FormValue(name)
}

func (ctx *contextDefault) Cookie(name string) (*Cookie, error) {
	c, err := ctx.request.Cookie(name)

	if err != nil {
		return nil, err
	}

	return &Cookie{
		Name:     c.Name,
		Value:    c.Value,
		Expires:  c.Expires,
		Domain:   c.Domain,
		Path:     c.Path,
		Secure:   c.Secure,
		HTTPOnly: c.HttpOnly,
	}, nil
}

func (ctx *contextDefault) SetCookie(c *Cookie) {
	http.SetCookie(ctx.writer, &http.Cookie{
		Name:    c.Name,
		Value:   c.Value,
		Expires: c.Expires,
	})
}

func (ctx *contextDefault) Extra(key string) interface{} {
	return ctx.extras[key]
}

func (ctx *contextDefault) AddExtra(key string, extra interface{}) {
	ctx.extras[key] = extra
}
