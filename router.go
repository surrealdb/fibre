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
	"net/url"
	"sort"
)

// Router stores routes used in request matching and handler dispatching.
type Router struct {
	fibre  *Fibre
	routes map[string][]*Route
}

// Route stores a handler for matching paths against requests.
type Route struct {
	Rank    int
	Path    string
	Method  string
	Handler HandlerFunc
}

// NewRouter returns a new Router instance.
func NewRouter(f *Fibre) *Router {
	return &Router{
		fibre:  f,
		routes: make(map[string][]*Route),
	}
}

// Add registers a new route with a matcher for the URL path.
func (r *Router) Add(meth, path string, hand HandlerFunc) {

	route := &Route{
		Path:    path,
		Method:  meth,
		Handler: hand,
	}

	// Rank the route
	route.Rank = route.rank()

	// Add the route to the tree
	r.routes[meth] = append(r.routes[meth], route)

	// Sort the routes according to rank
	sort.Sort(ranker(r.routes[meth]))

}

// Find dispatches the request to the handler whose path and method match
func (r *Router) Find(meth, path string, ctx *Context) (hand HandlerFunc) {

	routes := r.routes[meth]

	for _, route := range routes {
		if pa, ok := route.test(path); ok {
			ctx.param = pa
			return route.Handler
		}
	}

	// NOTE: This is slow, and also pointless if there is a catch all (*) route

	// allowed := make([]string, 0, len(r.routes))
	// for mech, routes := range r.routes {
	// 	if mech == meth {
	// 		continue
	// 	}
	// 	for _, route := range routes {
	// 		if _, ok := route.test(path); ok {
	// 			allowed = append(allowed, mech)
	// 		}
	// 	}
	// }

	// if len(allowed) >= 1 {
	// 	return func(c *Context) error {
	// 		return NewHTTPError(405)
	// 	}
	// }

	return func(c *Context) error {
		return NewHTTPError(404)
	}

}

func (r *Route) rank() (rank int) {

	for _, c := range r.Path {
		switch c {
		default:
			rank++
		case ':':
			rank += 100
		case '*':
			rank += 10000
		}
	}

	if rank > 20000 {
		panic("Path must only have 1 match all (*)")
	}

	if rank > 10000 && r.Path[len(r.Path)-1] != '*' {
		panic("Match all (*) must be at end of path")
	}

	return

}

func (r *Route) test(path string) (url.Values, bool) {

	var i int
	var j int
	var k string
	var v string

	param := make(url.Values)

	for i < len(r.Path) {

		switch {

		case r.Path[i] == '*':
			i++
			v, j = consumeCatch(j, path)
			param.Add("*", v)

		case r.Path[i] == ':':
			i++
			if k, i = consumeIdent(i, r.Path); len(k) == 0 {
				return nil, false
			}
			if v, j = consumeIdent(j, path); len(v) == 0 {
				return nil, false
			}
			param.Add(k, v)

		default:
			k, i = consumeChars(i, r.Path)
			v, j = consumeCount(j, path, len(k))
			if k != v {
				return nil, false
			}

		}

	}

	if hasRemaining(j, path) {
		return nil, false
	}

	return param, true

}

func hasRemaining(pos int, path string) bool {
	return len(path) > pos+1
}

func consumeCatch(pos int, path string) (s string, i int) {
	return path[pos:], len(path)
}

func consumeIdent(pos int, path string) (s string, i int) {
	for i = pos; i < len(path); i++ {
		if path[i] == '/' || path[i] == ':' || path[i] == '*' {
			break
		}
	}
	return path[pos:i], i
}

func consumeChars(pos int, path string) (s string, i int) {
	for i = pos; i < len(path); i++ {
		if path[i] == ':' || path[i] == '*' {
			break
		}
	}
	return path[pos:i], i
}

func consumeCount(pos int, path string, c int) (s string, i int) {
	if len(path) < pos+c {
		i = len(path)
	} else {
		i = pos + c
	}
	return path[pos:i], i
}
