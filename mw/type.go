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
	"github.com/abcum/fibre"
)

// TypeOpts defines options for the Type middleware.
type TypeOpts struct {
	AllowedContent map[string]bool
}

// Type defines middleware for checking the request content type.
func Type(opts ...*TypeOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			// Set defaults
			if len(opts) == 0 {
				opts = append(opts, &TypeOpts{})
			}

			// No config has been set
			if len(opts[0].AllowedContent) == 0 {
				return h(c)
			}

			if _, ok := opts[0].AllowedContent[c.Type()]; ok {
				return h(c)
			}

			return fibre.NewHTTPError(415)

		}
	}
}
