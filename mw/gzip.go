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

package mw

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/abcum/fibre"
)

var pool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

type zip struct {
	io.Writer
	http.ResponseWriter
}

func (z zip) Write(b []byte) (n int, err error) {
	return z.Writer.Write(b)
}

func (z zip) Flush() error {
	return z.Writer.(*gzip.Writer).Flush()
}

func (z zip) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return z.ResponseWriter.(http.Hijacker).Hijack()
}

func (z *zip) CloseNotify() <-chan bool {
	return z.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

// Gzip defines middleware for compressing response output.
func Gzip() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			// This is a websocket
			if c.Request().Header().Get("Upgrade") == "websocket" {
				return h(c)
			}

			// Set the accept-encoding header
			c.Response().Header().Add("Vary", "Accept-Encoding")

			// Check to see if the client can accept gzip encoding
			if strings.Contains(c.Request().Header().Get("Accept-Encoding"), "gzip") {

				// Get a zipper
				w := pool.Get().(*gzip.Writer)

				// Reset its io.Writer
				w.Reset(c.Response().Writer())

				defer func() {
					w.Close()
					pool.Put(w)
				}()

				// Specify the gzip encoding header
				c.Response().Header().Set("Content-Encoding", "gzip")

				// Set the response writer to the zipper
				c.Response().SetWriter(zip{
					Writer: w, ResponseWriter: c.Response().Writer(),
				})

			}

			if err := h(c); err != nil {
				c.Error(err)
			}

			return nil

		}
	}
}
