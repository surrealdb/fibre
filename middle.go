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
	"regexp"
)

func check(what string, checks ...string) bool {
	for _, check := range checks {
		if regex, err := regexp.Compile(check); err == nil {
			if regex.Match([]byte(what)) {
				return true
			}
		}
	}
	return false
}

// Host checks if the url host matches in the request. If the host matches the middleware will be invoked.
func (m MiddlewareFunc) Host(test ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if check(c.Request().URL().Host, test...) {
				return m(h)(c)
			}
			return h(c)
		}
	}
}

// Path checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) Path(test ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if check(c.Request().URL().Path, test...) {
				return m(h)(c)
			}
			return h(c)
		}
	}
}

// Scheme checks if the url scheme matches in the request. If the scheme matches the middleware will be invoked.
func (m MiddlewareFunc) Scheme(test ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if check(c.Request().URL().Scheme, test...) {
				return m(h)(c)
			}
			return h(c)
		}
	}
}

// UserAgent checks if the user agent matches in the request. If the user agent matches the middleware will be invoked.
func (m MiddlewareFunc) UserAgent(test ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			if check(c.Request().UserAgent(), test...) {
				return m(h)(c)
			}
			return h(c)
		}
	}
}
