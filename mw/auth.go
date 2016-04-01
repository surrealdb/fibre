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
	"bytes"

	"encoding/base64"

	"github.com/abcum/fibre"
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

			// Set defaults
			if len(opts) == 0 {
				opts = append(opts, &AuthOpts{})
			}

			// No config has been set
			if len(opts[0].User) == 0 && len(opts[0].Pass) == 0 {
				return h(c)
			}

			head := c.Request().Header().Get("Authorization")

			if head != "" && head[:5] == "Basic" {

				base, err := base64.StdEncoding.DecodeString(head[6:])

				if err == nil {

					cred := bytes.SplitN(base, []byte(":"), 2)

					if len(cred) == 2 && bytes.Equal(cred[0], opts[0].User) && bytes.Equal(cred[1], opts[0].Pass) {
						return h(c)
					}

				}

			}

			if opts[0].Realm != "" {
				c.Response().Header().Set("WWW-Authenticate", "Basic realm="+opts[0].Realm)
			}

			return fibre.NewHTTPError(401)

		}
	}
}
