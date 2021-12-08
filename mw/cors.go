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
	"net/http"
	"strconv"
	"strings"

	"github.com/surrealdb/fibre"
)

// CorsOpts defines options for the Cors middleware.
type CorsOpts struct {
	AllowedOrigin                 string
	AllowedMethods                []string
	AllowedHeaders                []string
	AccessControlMaxAge           int
	AccessControlAllowCredentials bool
}

var defaultCorsOpts = &CorsOpts{
	AllowedOrigin:                 "=",
	AllowedMethods:                []string{"GET", "PUT", "POST", "PATCH", "DELETE", "TRACE", "OPTIONS"},
	AllowedHeaders:                []string{"Accept", "Authorization", "Content-Type", "Origin"},
	AccessControlMaxAge:           600,
	AccessControlAllowCredentials: false,
}

// Cors defines middleware for setting and checking CORS headers,
func Cors(opts ...*CorsOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *CorsOpts

			switch len(opts) {
			case 0:
				config = defaultCorsOpts
			default:
				config = opts[0]
			}

			// This is a socket
			if c.IsSocket() {
				return h(c)
			}

			// No origin is set
			if !c.IsOrigin() {
				return h(c)
			}

			// Origin not allowed
			if config.AllowedOrigin != "*" && config.AllowedOrigin != "=" && config.AllowedOrigin != c.Origin() {
				return h(c)
			}

			// Normalize AllowedMethods and make comma-separated-values
			normedMethods := []string{}
			for _, allowedMethod := range config.AllowedMethods {
				normed := http.CanonicalHeaderKey(allowedMethod)
				normedMethods = append(normedMethods, normed)
			}

			// Normalize AllowedHeaders and make comma-separated-values
			normedHeaders := []string{}
			for _, allowedHeader := range config.AllowedHeaders {
				normed := http.CanonicalHeaderKey(allowedHeader)
				normedHeaders = append(normedHeaders, normed)
			}

			if len(normedMethods) > 0 {
				c.Response().Header().Set(fibre.HeaderAccessControlAllowMethods, strings.Join(normedMethods, ","))
			}

			if len(normedHeaders) > 0 {
				c.Response().Header().Set(fibre.HeaderAccessControlAllowHeaders, strings.Join(normedHeaders, ","))
			}

			if config.AccessControlMaxAge > 0 {
				c.Response().Header().Set(fibre.HeaderAccessControlMaxAge, strconv.Itoa(config.AccessControlMaxAge))
			}

			switch config.AllowedOrigin {
			default:
				c.Response().Header().Set(fibre.HeaderAccessControlAllowOrigin, config.AllowedOrigin)
			case "=":
				c.Response().Header().Set(fibre.HeaderAccessControlAllowOrigin, c.Origin())
			case "*":
				c.Response().Header().Set(fibre.HeaderAccessControlAllowOrigin, "*")
			}

			if config.AccessControlAllowCredentials {
				c.Response().Header().Set(fibre.HeaderAccessControlAllowCredentials, "true")
			}

			return h(c)

		}
	}
}
