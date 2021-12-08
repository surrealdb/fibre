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
	"bufio"
	"net"
	"net/http"
)

// Response wraps an http.Response
type Response struct {
	http.ResponseWriter
	fibre  *Fibre
	size   int64
	status int
	done   bool
}

// NewResponse creates a new instance of Response.
func NewResponse(i http.ResponseWriter, f *Fibre) *Response {
	return &Response{i, f, 0, 0, false}
}

// Size returns the current size, in bytes, of the response.
func (r *Response) Size() int64 {
	return r.size
}

// Done asserts whether or not the response has been sent.
func (r *Response) Done() bool {
	return r.done
}

// Writer returns the http.ResponseWriter instance for this Response.
func (r *Response) Writer() http.ResponseWriter {
	return r.ResponseWriter
}

// SetWriter sets the http.ResponseWriter instance for this Response.
func (r *Response) SetWriter(w http.ResponseWriter) {
	r.ResponseWriter = w
}

// Header returns the header map values for this Response.
func (r *Response) Header() http.Header {
	return r.ResponseWriter.Header()
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) WriteHeader(code int) {
	if r.done {
		return
	}
	r.status = code
	r.ResponseWriter.WriteHeader(code)
	r.done = true
}

// Write wraps and implements the http.Response.Write specification.
func (r *Response) Write(b []byte) (n int, err error) {
	n, err = r.ResponseWriter.Write(b)
	r.size += int64(n)
	return n, err
}

// Status returns the HTTP status code of the response.
func (r *Response) Status() int {
	return r.status
}

// Flush enables buffered data using http.Flusher.
func (r *Response) Flush() {
	r.ResponseWriter.(http.Flusher).Flush()
}

// Hijack enabled connection hijacking using http.Hijacker.
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.ResponseWriter.(http.Hijacker).Hijack()
}

// CloseNotify enables detecting when the underlying connection has gone away.
func (r *Response) CloseNotify() <-chan bool {
	return r.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (r *Response) reset(i http.ResponseWriter, f *Fibre) {
	r.fibre = f
	r.ResponseWriter = i
	r.done = false
	r.size = 0
	r.status = http.StatusOK
}
