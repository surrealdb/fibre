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

// UniqOpts defines options for the Uniq middleware.
type UniqOpts struct {
	HeaderKey string
}

var defaultUniqOpts = &UniqOpts{
	HeaderKey: "Request-Id",
}

// Uniq defines middleware for assigning a unique request id.
func Uniq(opts ...*UniqOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *UniqOpts

			switch len(opts) {
			case 0:
				config = defaultUniqOpts
			default:
				config = opts[0]
			}

			c.Response().Header().Set(config.HeaderKey, c.Uniq())

			return h(c)

		}
	}
}
