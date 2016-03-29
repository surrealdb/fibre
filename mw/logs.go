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

	"github.com/abcum/fibre"
	"github.com/labstack/gommon/color"
)

// Logs defines middleware for logging requests and responses.
func Logs() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) (err error) {

			code := "-"

			if err = h(c); err != nil {
				c.Error(err)
			}

			ip := c.IP()
			req := c.Request()
			res := c.Response()
			num := res.Status()
			now := c.Request().Start()

			met := req.Method
			url := req.URL.Path

			if c.Socket() != nil {
				met = "SOCK"
			}

			switch {
			case num >= 500:
				code = color.Red(num)
			case num >= 400:
				code = color.Yellow(num)
			case num >= 300:
				code = color.Cyan(num)
			case num >= 200:
				code = color.Blue(num)
			}

			c.Fibre().Logger().Infof("%s %s %s %s %s %d", color.Bold(code), ip, met, url, time.Since(now), res.Size())

			return nil

		}
	}
}
