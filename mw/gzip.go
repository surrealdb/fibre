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

package mw

import (
	"bufio"
	"compress/gzip"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/surrealdb/fibre"
)

var pool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

type zipper struct {
	gzip *gzip.Writer
	http.ResponseWriter
}

func (z *zipper) Setup() {

	// Get a gzip writer from the pool
	z.gzip = pool.Get().(*gzip.Writer)

	// Reset the pooled gzip writer
	z.gzip.Reset(z.ResponseWriter)

	// Remove any set length header
	z.ResponseWriter.Header().Del(fibre.HeaderContentLength)

	// Specify the gzip encoding header
	z.ResponseWriter.Header().Set(fibre.HeaderContentEncoding, "gzip")

}

func (z *zipper) Close() {
	if z.gzip != nil {
		z.gzip.Close()
		pool.Put(z.gzip)
	}
}

func (z *zipper) Write(b []byte) (n int, err error) {
	if z.gzip == nil {
		z.Setup()
	}
	if z.Header().Get(fibre.HeaderContentType) == "" {
		z.Header().Set(fibre.HeaderContentType, http.DetectContentType(b))
	}
	return z.gzip.Write(b)
}

func (z *zipper) WriteHeader(c int) {
	if z.gzip == nil {
		z.Setup()
	}
	if c == http.StatusNoContent {
		z.ResponseWriter.Header().Del(fibre.HeaderContentEncoding)
	}
	z.ResponseWriter.WriteHeader(c)
}

func (z *zipper) Flush() {
	if z.gzip != nil {
		z.gzip.Flush()
	}
	z.ResponseWriter.(http.Flusher).Flush()
}

func (z *zipper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return z.ResponseWriter.(http.Hijacker).Hijack()
}

func (z *zipper) CloseNotify() <-chan bool {
	return z.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Gzip defines middleware for compressing response output.
func Gzip() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			// This is a socket
			if c.IsSocket() {
				return h(c)
			}

			// Set the accept-encoding header

			c.Response().Header().Add(fibre.HeaderVary, "Accept-Encoding")

			// Check to see if the client can accept gzip encoding

			if strings.Contains(c.Request().Header().Get(fibre.HeaderAcceptEncoding), "gzip") {

				z := &zipper{ResponseWriter: c.Response().Writer()}

				c.Response().SetWriter(z)

				defer z.Close()

			}

			return h(c)

		}
	}
}
