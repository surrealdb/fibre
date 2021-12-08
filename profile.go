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
	"net/http/pprof"
)

func (f *Fibre) Pprof() {
	f.Get("/debug/pprof/", func(c *Context) error {
		pprof.Index(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/profile", func(c *Context) error {
		pprof.Profile(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/symbol", func(c *Context) error {
		pprof.Symbol(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/trace", func(c *Context) error {
		pprof.Trace(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/cmdline", func(c *Context) error {
		pprof.Cmdline(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/heap", func(c *Context) error {
		pprof.Handler("heap").ServeHTTP(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/block", func(c *Context) error {
		pprof.Handler("block").ServeHTTP(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/goroutine", func(c *Context) error {
		pprof.Handler("goroutine").ServeHTTP(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
	f.Get("/debug/pprof/threadcreate", func(c *Context) error {
		pprof.Handler("threadcreate").ServeHTTP(c.Response().ResponseWriter, c.Request().Request)
		return nil
	})
}
