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
	"time"

	"github.com/surrealdb/fibre"
)

// Logs defines middleware for logging requests and responses.
func Logs() fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) (err error) {

			err = h(c)

			ip := c.IP()
			req := c.Request()
			res := c.Response()
			num := res.Status()
			now := c.Request().Start()
			met := req.Method
			url := req.URL().Path
			log := c.Fibre().Logger().WithField("prefix", c.Fibre().Name())

			if c.Socket() != nil {
				met = "WS"
			}

			if err, ok := err.(*fibre.HTTPError); ok {
				num = err.Code()
				log = log.WithFields(err.Fields())
			}

			log = log.WithFields(map[string]interface{}{
				"ctx":     c,
				"id":      c.Uniq(),
				"ip":      ip,
				"url":     url,
				"size":    res.Size(),
				"status":  num,
				"method":  met,
				"latency": time.Since(now),
			})

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
			default:
				log.Info("Completed request")
			}

			return err

		}
	}
}
