// Copyright © 2016 Abcum Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fibre

import (
	"mime"
	"net"
	"os"
	// "path/filepath"
	"strings"

	// "io"
	"io/ioutil"

	"net/http"
	"net/url"

	"encoding/json"
	"encoding/xml"

	"github.com/gorilla/websocket"
	"gopkg.in/vmihailenco/msgpack.v2"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Context represents context for the current request.
type Context struct {
	fibre    *Fibre
	socket   *Socket
	request  *Request
	response *Response
	path     string
	param    url.Values
	query    url.Values
	store    map[string]interface{}
}

// NewContext creates a Context object.
func NewContext(req *Request, res *Response, f *Fibre) *Context {
	return &Context{
		fibre:    f,
		request:  req,
		response: res,
	}
}

// Fibre returns the fibre instance.
func (c *Context) Fibre() *Fibre {
	return c.fibre
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *Context) Error(err error) {
	c.fibre.errorHandler(err, c)
}

// Upgrade the http websocket connection.
func (c *Context) Upgrade() (err error) {
	if c.Request().Header().Get("Upgrade") == "websocket" {
		var sck *websocket.Conn
		req := c.request.Request
		res := c.response.ResponseWriter
		if sck, err = wsupgrader.Upgrade(res, req, nil); err == nil {
			c.socket = NewSocket(sck, c, c.Fibre())
		}
	}
	return
}

// Socket returns the websocket connection.
func (c *Context) Socket() *Socket {
	return c.socket
}

// Request returns the http request object.
func (c *Context) Request() *Request {
	return c.request
}

// Response returns the http response object.
func (c *Context) Response() *Response {
	return c.response
}

// Type returns the desired response mime type.
func (c *Context) Type() string {
	head := c.Request().Header().Get("Content-Type")
	cont, _, _ := mime.ParseMediaType(head)
	return cont
}

// Body returns the full content body.
func (c *Context) Body() []byte {
	body, _ := ioutil.ReadAll(c.Request().Body)
	return body
}

// Path returns the registered path for the handler.
func (c *Context) Path() string {
	return c.path
}

// Form returns form parameter by name.
func (c *Context) Form(name string) (v string) {
	return c.request.FormValue(name)
}

// Param returns path parameter by name.
func (c *Context) Param(name string) (v string) {
	return c.param.Get(name)
}

// Query returns query parameter by name.
func (c *Context) Query(name string) (v string) {
	return c.query.Get(name)
}

// Get retrieves data from the context.
func (c *Context) Get(key string) interface{} {
	if c.store == nil {
		return nil
	}
	return c.store[key]
}

// Set saves data in the context.
func (c *Context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = val
}

// Code sends a http response status code.
func (c *Context) Code(code int) (err error) {
	c.response.WriteHeader(code)
	return
}

// Text sends a text response with status code.
func (c *Context) Text(code int, data string) (err error) {
	c.response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.response.WriteHeader(code)
	c.response.Write([]byte(data))
	return
}

// HTML sends an html response with status code.
func (c *Context) HTML(code int, data string) (err error) {
	c.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.response.WriteHeader(code)
	c.response.Write([]byte(data))
	return
}

// XML sends a xml response with status code.
func (c *Context) XML(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "application/xml; charset=utf-8")
	c.response.WriteHeader(code)
	if !c.response.done {
		c.response.Write([]byte(xml.Header))
	}
	if data != nil {
		xml.NewEncoder(c.response.ResponseWriter).Encode(data)
	}
	return
}

// JSON sends a json response with status code.
func (c *Context) JSON(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.response.WriteHeader(code)
	if data != nil {
		json.NewEncoder(c.response.ResponseWriter).Encode(data)
	}
	return
}

// PACK sends a msgpack response with status code.
func (c *Context) PACK(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "application/msgpack; charset=utf-8")
	c.response.WriteHeader(code)
	if data != nil {
		msgpack.NewEncoder(c.response.ResponseWriter).Encode(data)
	}
	return
}

// Send sends the relevant response depending on the request type.
func (c *Context) Send(code int, data interface{}) (err error) {
	switch c.Type() {
	default:
		return c.JSON(code, data)
	case "application/xml":
		return c.XML(code, data)
	case "application/json":
		return c.JSON(code, data)
	case "application/msgpack":
		return c.PACK(code, data)
	}
}

// File sends a response with the content of a file.
func (c *Context) File(path string) (err error) {

	info, err := os.Stat(path)
	if err != nil {
		return NewHTTPError(404)
	}

	if info.IsDir() == false {
		file, err := os.Open(path)
		if err != nil {
			return NewHTTPError(404)
		}
		http.ServeContent(c.Response().ResponseWriter, c.Request().Request, info.Name(), info.ModTime(), file)
	}

	if info.IsDir() == true {
		file, err := os.Open(path + index)
		if err != nil {
			return NewHTTPError(404)
		}
		http.ServeContent(c.Response().ResponseWriter, c.Request().Request, info.Name(), info.ModTime(), file)
	}

	return nil

}

// Bind decodes the request body into the object.
func (c *Context) Bind(i interface{}) (err error) {
	switch c.Type() {
	case "application/xml":
		if err = xml.NewDecoder(c.Request().Body).Decode(i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	case "application/json":
		if err = json.NewDecoder(c.Request().Body).Decode(i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	case "application/msgpack":
		if err = msgpack.NewDecoder(c.Request().Body).Decode(i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	}
	return
}

// IP returns the ip address belonging to this context.
func (c *Context) IP() string {

	var i int
	var s string

	s = "X-Real-IP"
	for i = 0; i < 3; i++ {
		if i == 1 {
			s = strings.ToLower(s)
		}
		if i == 2 {
			s = strings.ToUpper(s)
		}
		if ip := c.Request().Header().Get(s); ip != "" {
			return ip
		}
	}

	s = "X-Forwarded-For"
	for i = 0; i < 3; i++ {
		if i == 1 {
			s = strings.ToLower(s)
		}
		if i == 2 {
			s = strings.ToUpper(s)
		}
		if ip := c.Request().Header().Get(s); ip != "" {
			return strings.Split(ip, ", ")[0]
		}
	}

	addr := c.Request().RemoteAddr
	addr, _, _ = net.SplitHostPort(addr)
	return addr

}

func (c *Context) reset(r *http.Request, w http.ResponseWriter, f *Fibre) {

	// Set the fibre instance
	c.fibre = f

	// Reset the query and store vars
	c.param = nil
	c.query = nil
	c.store = nil

	// Reset the request and response
	c.socket = nil
	c.request.reset(r, f)
	c.response.reset(w, f)

	// Reset the url query paramaters
	c.query = r.URL.Query()

}
