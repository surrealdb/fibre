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
)

// Logs defines middleware for logging requests and responses.
func Logs() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) (err error) {

			if err = h(c); err != nil {
				c.Error(err)
			}

			ip := c.IP()
			req := c.Request()
			res := c.Response()
			num := res.Status()
			now := c.Request().Start()

			met := req.Method
			url := req.URL().Path

			if c.Socket() != nil {
				met = "SOCK"
			}

			log := c.Fibre().Logger().WithFields(map[string]interface{}{
				"prefix": c.Fibre().Name(),
				"ip":     ip,
				"url":    url,
				"size":   res.Size(),
				"status": num,
				"method": met,
				"time":   time.Since(now),
			})

			if id := c.Get("id"); id != nil {
				log = log.WithField("id", id)
			}

			if err != nil {
				log = log.WithError(err)
			}

			switch {
			case num >= 500:
				log.Error("Completed request")
			case num >= 400:
				log.Warn("Completed request")
			case num >= 300:
				log.Info("Completed request")
			case num >= 200:
				log.Info("Completed request")
			}

			return nil

		}
	}
}
