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
	"os"
)

// Extra keys
const (
	ExtraUser = "github.com/mbrlabs/zion/security:extra_user"
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

// NewContext creates a new context
func NewContext(z *Zion, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		zion:      z,
		writer:    w,
		request:   r,
		urlParams: make(map[string]string),
		extras:    make(map[string]interface{}),
	}
}

// Render renders a template given the template name and some data
func (ctx *Context) Render(name string, data interface{}) {
	ctx.zion.config.TemplateEngine.Render(name, data, ctx.writer)
}

// String simply sends the given string to the client
func (ctx *Context) String(str string) {
	fmt.Fprintf(ctx.writer, "%s", str)
}

// Json coverts the data to a json, sets the content type & sends it to the client
func (ctx *Context) Json(data interface{}) {
	data, err := json.Marshal(data)
	if err != nil {
		fmt.Println("[ERROR] " + err.Error())
	}
	ctx.writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(ctx.writer, "%s", data)
}

// RawJson sets the content type & sends the json to the client
func (ctx *Context) RawJson(json string) {
	ctx.writer.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(ctx.writer, "%s", json)
}

// File sends a file to the client. This method implements the If-Modified-Since/Last-Modified headers
// and sends files only if the client has an older version of the file.
func (ctx *Context) File(path string) {
	// open file
	file, err := os.Open(path)
	defer file.Close()

	// return if file not found
	if err != nil {
		ctx.SendStatus(http.StatusNotFound)
		return
	}

	// return if file info cound not be read
	fileInfo, err := file.Stat()
	if err != nil {
		ctx.SendStatus(http.StatusNotFound)
		return
	}

	// serve file. this also considers Last-Modified/If-Modified-Since
	http.ServeContent(ctx.writer, ctx.request, fileInfo.Name(), fileInfo.ModTime(), file)
}

// SendStatus replies only with a status response
func (ctx *Context) SendStatus(status int) {
	ctx.writer.WriteHeader(status)
}

// Redirect redirects the user to the given relative url
func (ctx *Context) Redirect(path string) {
	http.Redirect(ctx.writer, ctx.request, path, http.StatusTemporaryRedirect)
}

// Method returns the http method of th request
func (ctx *Context) Method() string {
	return ctx.request.Method
}

// RequestHeader returns the http header of the request
func (ctx *Context) RequestHeader() http.Header {
	return ctx.request.Header
}

// ResponseHeader returns the http header of the response
func (ctx *Context) ResponseHeader() http.Header {
	return ctx.writer.Header()
}

// Path returns path of the request
func (ctx *Context) Path() string {
	return ctx.request.URL.Path
}

// URLParams returnes the url parameters (/home/:param/*file), parsed by the router
func (ctx *Context) URLParams() map[string]string {
	return ctx.urlParams
}

// URLParam returns one url parameter by name
func (ctx *Context) URLParam(name string) string {
	return ctx.urlParams[name]
}

// FormValue returns one url encoded post parameter by name
func (ctx *Context) FormValue(name string) string {
	return ctx.request.FormValue(name)
}

// Cookie gets a cookie in the request by name
func (ctx *Context) Cookie(name string) (*http.Cookie, error) {
	return ctx.request.Cookie(name)
}

// SetCookie sets a cookie for the response
func (ctx *Context) SetCookie(c *http.Cookie) {
	http.SetCookie(ctx.writer, c)
}

// Extra returns an user defined extra
func (ctx *Context) Extra(key string) interface{} {
	return ctx.extras[key]
}

// AddExtra adds a user defined extra
func (ctx *Context) AddExtra(key string, extra interface{}) {
	ctx.extras[key] = extra
}
