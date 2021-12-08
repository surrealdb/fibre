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
	"time"

	"context"

	"github.com/surrealdb/fibre"
)

// QuitOpts defines options for the Sign middleware.
type QuitOpts struct {
	Timeout time.Duration
}

// Quit defines middleware for timing out a connection.
func Quit(opts ...*QuitOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) (err error) {

			var config *QuitOpts

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
			if config.Timeout == 0 {
				return h(c)
			}

			ctx, cancel := context.WithTimeout(c.Context(), config.Timeout)

			c = c.WithContext(ctx)

			defer cancel()

			return h(c)

		}
	}
}
