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
	"time"

	"context"

	"github.com/abcum/fibre"
)

// QuitOpts defines options for the Sign middleware.
type QuitOpts struct {
	Timeout time.Duration
}

// Quit defines middleware for timing out a connection.
func Quit(opts ...*QuitOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) (err error) {

			// Set defaults
			if len(opts) == 0 {
				return h(c)
			}

			// No config has been set
			if opts[0].Timeout == 0 {
				return h(c)
			}

			// This is a websocket
			if c.Request().Header().Get("Upgrade") == "websocket" {
				return h(c)
			}

			ctx, cancel := context.WithTimeout(c.Context(), opts[0].Timeout)

			c = c.WithContext(ctx)

			defer cancel()

			return h(c)

		}
	}
}
