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

// NewrelicOpts defines options for the Info middleware.
type NewrelicOpts struct {
	Name    string
	License string
}

// Newrelic returns a middleware function for newrelic monitoring.
func Newrelic(opts ...*NewrelicOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			// Set defaults
			if len(opts) == 0 {
				opts = append(opts, &NewrelicOpts{})
			}

			// No config has been set
			if len(opts[0].Name) == 0 || len(opts[0].License) == 0 {
				return h(c)
			}

			if agent == nil {
				config := newrelic.NewConfig(opts[0].Name, opts[0].License)
				agent, _ = newrelic.NewApplication(config)
			}

			txn := agent.StartTransaction(c.Request().URL().Path, c.Response().ResponseWriter, c.Request().Request)

			defer txn.End()

			return h(c)

		}
	}
}
