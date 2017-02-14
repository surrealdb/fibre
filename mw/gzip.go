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

type zipper struct {
	io.Writer
	http.ResponseWriter
}

func (w zipper) Write(b []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w zipper) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w zipper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func (w *zipper) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

var writerPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

// Gzip defines middleware for compressing response output.
func Gzip() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			// This is a websocket
			if c.Request().Header().Get("Upgrade") == "websocket" {
				return h(c)
			}

			c.Response().Header().Add("Vary", "Accept-Encoding")

			if strings.Contains(c.Request().Header().Get("Accept-Encoding"), "gzip") {
				w := writerPool.Get().(*gzip.Writer)
				w.Reset(c.Response().Writer())
				defer func() {
					w.Close()
					writerPool.Put(w)
				}()
				gw := zipper{Writer: w, ResponseWriter: c.Response().Writer()}
				c.Response().Header().Set("Content-Encoding", "gzip")
				c.Response().SetWriter(gw)
			}

			if err := h(c); err != nil {
				c.Error(err)
			}

			return nil

		}
	}
}
