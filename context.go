// Copyright © SurrealDB Ltd
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
	"context"
	"mime"
	"net"
	"os"
	"strings"
	"time"

	"io/ioutil"

	"net/http"
	"net/url"

	"encoding/xml"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/ksuid"
	"github.com/ugorji/go/codec"
)

// Context represents context for the current request.
type Context struct {
	fibre    *Fibre
	socket   *Socket
	request  *Request
	response *Response
	ctx      context.Context
	uniq     string
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

func (c *Context) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *Context) WithContext(ctx context.Context) *Context {
	n := new(Context)
	*n = *c
	n.ctx = ctx
	return n
}

// Fibre returns the fibre instance.
func (c *Context) Fibre() *Fibre {
	return c.fibre
}

func (c *Context) Uniq() string {
	return c.uniq
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *Context) Error(err error) {
	c.fibre.errorHandler(err, c)
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

// IsTLS returns true if the request was made over TLS.
func (c *Context) IsTLS() bool {
	return c.Request().TLS != nil
}

// IsOrigin returns true if the request specifies an origin.
func (c *Context) IsOrigin() bool {
	return c.Request().Header().Get(HeaderOrigin) != ""
}

// IsSocket returns true if the request is made over WebSocket.
func (c *Context) IsSocket() bool {
	return c.Request().Header().Get(HeaderUpgrade) == "websocket"
}

// IsComplete returns true if the response has been closed.
func (c *Context) IsComplete() bool {
	return c.Response().Done()
}

// Type returns the desired response mime type.
func (c *Context) Type() string {
	head := c.Request().Header().Get("Content-Type")
	cont, _, _ := mime.ParseMediaType(head)
	return cont
}

// Head returns the processed headers.
func (c *Context) Head() map[string]string {
	head := map[string]string{}
	for k := range c.Request().Header() {
		head[k] = c.Request().Header().Get(k)
	}
	return head
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

// Origin returns the request origin if specified.
func (c *Context) Origin() (v string) {
	return c.Request().Header().Get(HeaderOrigin)
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
	if v, ok := c.store[key]; ok {
		return v
	}
	return nil
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

// Data sends a response with status code and mime type.
func (c *Context) Data(code int, data interface{}, mime string) (err error) {
	c.response.Header().Set("Content-Type", mime)
	c.response.WriteHeader(code)
	switch conv := data.(type) {
	case []byte:
		c.response.Write(conv)
	case string:
		c.response.Write([]byte(conv))
	}
	return
}

// Text sends a text response with status code.
func (c *Context) Text(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.response.WriteHeader(code)
	switch conv := data.(type) {
	case []byte:
		c.response.Write(conv)
	case string:
		c.response.Write([]byte(conv))
	}
	return
}

// HTML sends an html response with status code.
func (c *Context) HTML(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.response.WriteHeader(code)
	switch conv := data.(type) {
	case []byte:
		c.response.Write(conv)
	case string:
		c.response.Write([]byte(conv))
	}
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
		return xml.NewEncoder(c.response).Encode(data)
	}
	return
}

// JSON sends a json response with status code.
func (c *Context) JSON(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.response.WriteHeader(code)
	if data != nil {
		return codec.NewEncoder(c.response, &jh).Encode(data)
	}
	return
}

// CBOR sends a cbor response with status code.
func (c *Context) CBOR(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "application/cbor; charset=utf-8")
	c.response.WriteHeader(code)
	if data != nil {
		return codec.NewEncoder(c.response, &ch).Encode(data)
	}
	return
}

// PACK sends a msgpack response with status code.
func (c *Context) PACK(code int, data interface{}) (err error) {
	c.response.Header().Set("Content-Type", "application/msgpack; charset=utf-8")
	c.response.WriteHeader(code)
	if data != nil {
		return codec.NewEncoder(c.response, &mh).Encode(data)
	}
	return
}

// Send sends the relevant response depending on the request type.
func (c *Context) Send(code int, data interface{}) (err error) {
	switch c.Type() {
	default:
		return c.Text(code, data)
	case "application/xml":
		return c.XML(code, data)
	case "application/json":
		return c.JSON(code, data)
	case "application/cbor":
		return c.CBOR(code, data)
	case "application/msgpack":
		return c.PACK(code, data)
	case "application/vnd.api+json":
		return c.JSON(code, data)
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
		http.ServeContent(c.response, c.Request().Request, info.Name(), info.ModTime(), file)
	}

	if info.IsDir() == true {
		file, err := os.Open(path + index)
		if err != nil {
			return NewHTTPError(404)
		}
		http.ServeContent(c.response, c.Request().Request, info.Name(), info.ModTime(), file)
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
		if err = codec.NewDecoder(c.Request().Body, &jh).Decode(i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	case "application/cbor":
		if err = codec.NewDecoder(c.Request().Body, &ch).Decode(i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	case "application/msgpack":
		if err = codec.NewDecoder(c.Request().Body, &mh).Decode(i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	case "application/x-www-form-urlencoded":
		obj := map[string]interface{}{}
		if err = c.request.ParseForm(); err != nil {
			err = NewHTTPError(400, err.Error())
		}
		for k, v := range c.request.Form {
			if len(v) == 1 {
				obj[k] = v[0]
			} else {
				obj[k] = v
			}
		}
		if err = mapstructure.Decode(obj, i); err != nil {
			err = NewHTTPError(400, err.Error())
		}
	}
	return
}

// IP returns the ip address belonging to this context.
func (c *Context) IP() net.IP {

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
			return net.ParseIP(ip)
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
			return net.ParseIP(strings.Split(ip, ", ")[0])
		}
	}

	addr := c.Request().RemoteAddr
	addr, _, _ = net.SplitHostPort(addr)
	return net.ParseIP(addr)

}

// Redirect redirects the http request to a different url.
func (c *Context) Redirect(code int, url string) (err error) {
	c.Response().Header().Set(HeaderLocation, url)
	c.Response().WriteHeader(code)
	return nil
}

// Upgrade upgrades the http request to a websocket connection.
func (c *Context) Upgrade(protocols ...string) (err error) {

	wes := websocket.Upgrader{
		ReadBufferSize:    4096,
		WriteBufferSize:   4096,
		EnableCompression: false,
		Subprotocols:      protocols,
		HandshakeTimeout:  time.Second * 10,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	if websocket.IsWebSocketUpgrade(c.Request().Request) {

		var sck *websocket.Conn
		res := c.response
		req := c.request.Request
		pro := websocket.Subprotocols(c.Request().Request)

		if len(protocols) > 0 && !in(protocols, pro) {
			return NewHTTPError(415, "Unsupported Media Type")
		}

		if sck, err = wes.Upgrade(res, req, res.Header()); err != nil {
			return NewHTTPError(426, "Upgrade required")
		}

		c.socket = NewSocket(sck, c, c.Fibre())

		return nil

	}

	return NewHTTPError(426, "Upgrade required")

}

func (c *Context) reset(r *http.Request, w http.ResponseWriter, f *Fibre) {

	// Set the fibre instance
	c.fibre = f

	// Set an id for this connection
	c.uniq = ksuid.New().String()

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
