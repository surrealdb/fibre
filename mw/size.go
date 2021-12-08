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
	"github.com/surrealdb/fibre"
)

// SizeOpts defines options for the Head middleware.
type SizeOpts struct {
	AllowedLength int64
}

// Size defines middleware for checking the request content length.
func Size(opts ...*SizeOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *SizeOpts

			switch len(opts) {
			case 0:
				return h(c)
			default:
				config = opts[0]
			}

			// This is a socket
			if c.IsSocket() {
				return h(c)
			}

			// No config has been set
			if config.AllowedLength == 0 {
				return h(c)
			}

			// Content length is within allowed limits
			if c.Request().ContentLength <= config.AllowedLength {
				return h(c)
			}

			return fibre.NewHTTPError(413)

		}
	}
}
