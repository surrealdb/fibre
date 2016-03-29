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
	"runtime"

	"github.com/abcum/fibre"
)

// Fail defines middleware for recovering from panics,
func Fail() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			defer func() {
				if err := recover(); err != nil {
					trace := make([]byte, 1<<16)
					n := runtime.Stack(trace, false)
					c.Fibre().Logger().Errorf("%v\n stack trace %s", err, trace[:n])
					c.Error(fmt.Errorf("panic recover\n %v\n stack trace %d bytes\n %s", err, n, trace[:n]))
				}
			}()

			return h(c)

		}
	}
}
