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
	"strings"
)

// ----------------------------------------------------------------------------------------------------

// HostIs checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) HostIs(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().URL().Host == test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// HostMatches checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) HostMatches(tests ...regexp.Regexp) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if test.MatchString(c.Request().URL().Host) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// HostBegsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) HostBegsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasPrefix(c.Request().URL().Host, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// HostEndsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) HostEndsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasSuffix(c.Request().URL().Host, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// ----------------------------------------------------------------------------------------------------

// PathIs checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) PathIs(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().URL().Path == test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// PathMatches checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) PathMatches(tests ...regexp.Regexp) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if test.MatchString(c.Request().URL().Path) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// PathBegsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) PathBegsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasPrefix(c.Request().URL().Path, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// PathEndsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) PathEndsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasSuffix(c.Request().URL().Path, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// ----------------------------------------------------------------------------------------------------

// SchemeIs checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) SchemeIs(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().URL().Scheme == test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// SchemeMatches checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) SchemeMatches(tests ...regexp.Regexp) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if test.MatchString(c.Request().URL().Scheme) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// SchemeBegsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) SchemeBegsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasPrefix(c.Request().URL().Scheme, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// SchemeEndsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) SchemeEndsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasSuffix(c.Request().URL().Scheme, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// ----------------------------------------------------------------------------------------------------

// AgentIs checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) AgentIs(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().UserAgent() == test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// AgentMatches checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) AgentMatches(tests ...regexp.Regexp) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if test.MatchString(c.Request().UserAgent()) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// AgentBegsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) AgentBegsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasPrefix(c.Request().UserAgent(), test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// AgentEndsWith checks if the url path matches in the request. If the path matches the middleware will be invoked.
func (m MiddlewareFunc) AgentEndsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasSuffix(c.Request().UserAgent(), test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}
