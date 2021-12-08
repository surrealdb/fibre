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
	"net/http"
	"sync"
	"time"
)

type (
	// Fibre represents an HTTP server
	Fibre struct {
		pool         sync.Pool
		name         string
		wait         time.Duration
		itimeout     time.Duration
		rtimeout     time.Duration
		wtimeout     time.Duration
		logger       *Logger
		router       *Router
		middleware   Middleware
		errorHandler HTTPErrorHandler
	}

	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error, *Context)

	// Middleware stores loaded middleware
	Middleware []MiddlewareFunc

	// HandlerFunc represents a request handler
	HandlerFunc func(*Context) error

	// MiddlewareFunc represents a request middleware
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

const (
	// HEAD ...
	HEAD = "HEAD"
	// GET ...
	GET = "GET"
	// PUT ...
	PUT = "PUT"
	// POST ...
	POST = "POST"
	// PATCH ...
	PATCH = "PATCH"
	// TRACE ...
	TRACE = "TRACE"
	// DELETE ...
	DELETE = "DELETE"
	// OPTIONS ...
	OPTIONS = "OPTIONS"
	// CONNECT ...
	CONNECT = "CONNECT"
)

const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthenticate        = "WWW-Authenticate"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderOrigin              = "Origin"
	HeaderServer              = "Server"
	HeaderSetCookie           = "Set-Cookie"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXPoweredBy          = "X-Powered-By"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXUrlScheme          = "X-Url-Scheme"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderPublicKeyPins           = "Public-Key-Pins"
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXCSRFToken              = "X-CSRF-Token"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
)

var (
	index = "index.html"

	methods = [...]string{HEAD, GET, PUT, POST, PATCH, TRACE, DELETE, OPTIONS, CONNECT}
)

// Server creates a new server instance.
func Server() (f *Fibre) {

	f = &Fibre{}

	// Set the default name
	f.name = "fibre"

	// Setup a new logger
	f.logger = NewLogger(f)

	// Setup a new router
	f.router = NewRouter(f)

	// Setup the default error handler
	f.SetHTTPErrorHandler(f.defaultErrorHandler)

	// Setup a new context pool
	f.pool.New = func() interface{} {
		return NewContext(new(Request), new(Response), f)
	}

	return

}

// Name returns the instance name.
func (f *Fibre) Name() string {
	return f.name
}

// Logger returns the logger instance.
func (f *Fibre) Logger() *Logger {
	return f.logger
}

// Router returns the router instance.
func (f *Fibre) Router() *Router {
	return f.router
}

// SetName sets the instance name.
func (f *Fibre) SetName(name string) {
	f.name = name
}

// SetLogLevel sets the logger log level.
func (f *Fibre) SetLogLevel(l string) {
	f.Logger().SetLevel(l)
}

// SetLogFormat sets the logger log format.
func (f *Fibre) SetLogFormat(l string) {
	f.Logger().SetFormat(l)
}

// SetIdleTimeout sets the max idle time for a keepalive connection.
func (f *Fibre) SetIdleTimeout(wait string) {
	f.itimeout, _ = time.ParseDuration(wait)
}

// SetReadTimeout sets the max duration for reading requests.
func (f *Fibre) SetReadTimeout(wait string) {
	f.rtimeout, _ = time.ParseDuration(wait)
}

// SetWriteTimeout sets the max duration for writing responses.
func (f *Fibre) SetWriteTimeout(wait string) {
	f.wtimeout, _ = time.ParseDuration(wait)
}

// SetHTTPErrorHandler registers a custom Echo.HTTPErrorHandler.
func (f *Fibre) SetHTTPErrorHandler(h HTTPErrorHandler) {
	f.errorHandler = h
}

// Use adds a middleware function
func (f *Fibre) Use(m MiddlewareFunc) MiddlewareFunc {
	f.middleware = append(f.middleware, m)
	return m
}

// Head adds a HEAD route > handler to the router.
func (f *Fibre) Head(p string, h HandlerFunc) {
	f.router.Add(HEAD, p, h)
}

// Get adds a GET route > handler to the router.
func (f *Fibre) Get(p string, h HandlerFunc) {
	f.router.Add(GET, p, h)
}

// Put adds a PUT route > handler to the router.
func (f *Fibre) Put(p string, h HandlerFunc) {
	f.router.Add(PUT, p, h)
}

// Post adds a POST route > handler to the router.
func (f *Fibre) Post(p string, h HandlerFunc) {
	f.router.Add(POST, p, h)
}

// Patch adds a PATCH route > handler to the router.
func (f *Fibre) Patch(p string, h HandlerFunc) {
	f.router.Add(PATCH, p, h)
}

// Trace adds a TRACE route > handler to the router.
func (f *Fibre) Trace(p string, h HandlerFunc) {
	f.router.Add(TRACE, p, h)
}

// Delete adds a DELETE route > handler to the router.
func (f *Fibre) Delete(p string, h HandlerFunc) {
	f.router.Add(DELETE, p, h)
}

// Options adds an OPTIONS route > handler to the router.
func (f *Fibre) Options(p string, h HandlerFunc) {
	f.router.Add(OPTIONS, p, h)
}

// Connect adds a CONNECT route > handler to the router.
func (f *Fibre) Connect(p string, h HandlerFunc) {
	f.router.Add(CONNECT, p, h)
}

// Any adds a route > handler to the router for all HTTP methods.
func (f *Fibre) Any(p string, h HandlerFunc) {
	for _, m := range methods {
		f.router.Add(m, p, h)
	}
}

// Dir serves a folder.
func (f *Fibre) Dir(p, dir string) {
	f.Get(p+"*", func(c *Context) error {
		return c.File(dir + c.Param("*"))
	})
}

// File serves a file.
func (f *Fibre) File(p, file string) {
	f.Get(p, func(c *Context) error {
		return c.File(file)
	})
}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (f *Fibre) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := f.pool.Get().(*Context)

	c.reset(r, w, f)

	p := f.router.Find(r.Method, r.URL.Path, c)

	// Catch all errors before sending any output
	h := func(c *Context) (err error) {
		if err = p(c); err != nil {
			c.Error(err)
		}
		return
	}

	// Chain middleware with handler in the end
	for i := len(f.middleware) - 1; i >= 0; i-- {
		h = f.middleware[i](h)
	}

	// Execute middleware and request chain
	if err := h(c); err != nil {
		c.Error(err)
	}

	f.pool.Put(c)

}
