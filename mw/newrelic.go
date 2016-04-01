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
	"github.com/yvasiyarov/go-metrics"
	"github.com/yvasiyarov/gorelic"
)

var agent *gorelic.Agent

// NewrelicOpts defines options for the Info middleware.
type NewrelicOpts struct {
	Name    []byte
	License []byte
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
				agent = gorelic.NewAgent()
				agent.Verbose = false
				agent.CollectHTTPStat = true
				agent.CollectHTTPStatuses = true
				agent.HTTPTimer = metrics.NewTimer()
				agent.HTTPStatusCounters = make(map[int]metrics.Counter)
				agent.NewrelicName = string(opts[0].Name)
				agent.NewrelicLicense = string(opts[0].License)
				agent.NewrelicPollInterval = 60
				agent.Run()
			}

			if err := h(c); err != nil {
				c.Error(err)
			}

			agent.HTTPTimer.UpdateSince(c.Request().Start())
			agent.HTTPStatusCounters[c.Response().Status()].Inc(1)

			return nil

		}
	}
}
