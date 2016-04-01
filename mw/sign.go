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
	"fmt"

	"github.com/abcum/fibre"
	"github.com/dgrijalva/jwt-go"
)

// SignOpts defines options for the Sign middleware.
type SignOpts struct {
	Key []byte
	Fnc func(*fibre.Context, *jwt.Token) error
}

// Sign defines middleware for JWT authentication.
func Sign(opts ...*SignOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			// Set defaults
			if len(opts) == 0 {
				return h(c)
			}

			// No config has been set
			if len(opts[0].Key) == 0 {
				return h(c)
			}

			head := c.Request().Header().Get("Authorization")

			if head != "" && head[:6] == "Bearer" {

				token, err := jwt.Parse(head[7:], func(token *jwt.Token) (interface{}, error) {

					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					}

					return opts[0].Key, nil

				})

				if err == nil && token.Valid {
					if opts[0].Fnc != nil {
						if err := opts[0].Fnc(c, token); err != nil {
							return fibre.NewHTTPError(401)
						}
					}
					return h(c)
				}

			}

			return fibre.NewHTTPError(401)

		}
	}
}
