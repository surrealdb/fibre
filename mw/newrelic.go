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
	"github.com/newrelic/go-agent"
)

var agent newrelic.Application

// NewrelicOpts defines options for the Newrelic middleware.
type NewrelicOpts struct {
	Name    string
	License string
}

// Newrelic returns a middleware function for newrelic monitoring.
func Newrelic(opts ...*NewrelicOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *NewrelicOpts

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
			if len(config.Name) == 0 || len(config.License) == 0 {
				return h(c)
			}

			if agent == nil {
				config := newrelic.NewConfig(config.Name, config.License)
				agent, _ = newrelic.NewApplication(config)
			}

			txn := agent.StartTransaction(c.Request().URL().Path, c.Response(), c.Request().Request)

			defer txn.End()

			return h(c)

		}
	}
}
