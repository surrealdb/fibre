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

// InfoOpts defines options for the Info middleware.
type InfoOpts struct {
	PoweredBy string
}

var defaultInfoOpts = &InfoOpts{
	PoweredBy: "Fibre",
}

// Info defines middleware for specifying the server powered-by header.
func Info(opts ...*InfoOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *InfoOpts

			switch len(opts) {
			case 0:
				config = defaultInfoOpts
			default:
				config = opts[0]
			}

			c.Response().Header().Set(fibre.HeaderXPoweredBy, config.PoweredBy)

			return h(c)

		}
	}
}
