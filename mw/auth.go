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
	"bytes"

	"encoding/base64"

	"github.com/surrealdb/fibre"
)

// AuthOpts defines options for the Auth middleware.
type AuthOpts struct {
	User  []byte
	Pass  []byte
	Realm string
}

// Auth defines middleware for basic http authentication.
func Auth(opts ...*AuthOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *AuthOpts

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
			if len(config.User) == 0 && len(config.Pass) == 0 {
				return h(c)
			}

			head := c.Request().Header().Get(fibre.HeaderAuthorization)

			if head != "" && head[:5] == "Basic" {

				base, err := base64.StdEncoding.DecodeString(head[6:])

				if err == nil {

					cred := bytes.SplitN(base, []byte(":"), 2)

					if len(cred) == 2 && bytes.Equal(cred[0], config.User) && bytes.Equal(cred[1], config.Pass) {
						return h(c)
					}

				}

			}

			if config.Realm != "" {
				c.Response().Header().Set(fibre.HeaderAuthenticate, "Basic realm="+config.Realm)
			}

			return fibre.NewHTTPError(401)

		}
	}
}
