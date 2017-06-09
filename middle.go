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

var MiddlewareSkip = func(h HandlerFunc) HandlerFunc {
	return func(c *Context) error {
		return h(c)
	}
}

// ----------------------------------------------------------------------

// HostIs checks if the request host is exactly equal to
// a value, and if it is then the middleware will be invoked.
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

// HostIsNot checks if the request host is not exactly equal to
// a value, and if it isn't then the middleware will be invoked.
func (m MiddlewareFunc) HostIsNot(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().URL().Host != test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// HostMatches checks if the request host matches a regular
// expression, and if it does then the middleware will be invoked.
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

// HostBegsWith checks if the request host begins with a value,
// and if it does then the middleware will be invoked.
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

// HostEndsWith checks if the request host ends with a value,
// and if it does then the middleware will be invoked.
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

// ----------------------------------------------------------------------

// PathIs checks if the request oath is exactly equal to
// a value, and if it is then the middleware will be invoked.
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

// PathIsNot checks if the request oath is not exactly equal to
// a value, and if it isn't then the middleware will be invoked.
func (m MiddlewareFunc) PathIsNot(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().URL().Path != test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// PathMatches checks if the request oath matches a regular
// expression, and if it does then the middleware will be invoked.
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

// PathBegsWith checks if the request oath begins with a value,
// and if it does then the middleware will be invoked.
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

// PathEndsWith checks if the request oath ends with a value,
// and if it does then the middleware will be invoked.
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

// ----------------------------------------------------------------------

// SchemeIs checks if the request scheme is exactly equal to
// a value, and if it is then the middleware will be invoked.
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

// SchemeIsNot checks if the request scheme is not exactly equal to
// a value, and if it isn't then the middleware will be invoked.
func (m MiddlewareFunc) SchemeIsNot(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().URL().Scheme != test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// SchemeMatches checks if the request scheme matches a regular
// expression, and if it does then the middleware will be invoked.
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

// SchemeBegsWith checks if the request scheme begins with a value,
// and if it does then the middleware will be invoked.
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

// SchemeEndsWith checks if the request scheme ends with a value,
// and if it does then the middleware will be invoked.
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

// ----------------------------------------------------------------------

// AgentIs checks if the request agent is exactly equal to
// a value, and if it is then the middleware will be invoked.
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

// AgentIsNot checks if the request agent is not exactly equal to
// a value, and if it isn't then the middleware will be invoked.
func (m MiddlewareFunc) AgentIsNot(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().UserAgent() != test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// AgentMatches checks if the request agent matches a regular
// expression, and if it does then the middleware will be invoked.
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

// AgentBegsWith checks if the request agent begins with a value,
// and if it does then the middleware will be invoked.
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

// AgentEndsWith checks if the request agent ends with a value,
// and if it does then the middleware will be invoked.
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

// ----------------------------------------------------------------------

// MethodIs checks if the request method is exactly equal to
// a value, and if it is then the middleware will be invoked.
func (m MiddlewareFunc) MethodIs(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().Request.Method == test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// MethodIsNot checks if the request method is not exactly equal to
// a value, and if it isn't then the middleware will be invoked.
func (m MiddlewareFunc) MethodIsNot(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if c.Request().Request.Method != test {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// MethodMatches checks if the request method matches a regular
// expression, and if it does then the middleware will be invoked.
func (m MiddlewareFunc) MethodMatches(tests ...regexp.Regexp) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if test.MatchString(c.Request().Request.Method) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// MethodBegsWith checks if the request method begins with a value,
// and if it does then the middleware will be invoked.
func (m MiddlewareFunc) MethodBegsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasPrefix(c.Request().Request.Method, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}

// MethodEndsWith checks if the request method ends with a value,
// and if it does then the middleware will be invoked.
func (m MiddlewareFunc) MethodEndsWith(tests ...string) MiddlewareFunc {
	return func(h HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			for _, test := range tests {
				if strings.HasSuffix(c.Request().Request.Method, test) {
					return m(h)(c)
				}
			}
			return h(c)
		}
	}
}
