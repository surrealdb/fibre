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
	"fmt"

	"github.com/surrealdb/fibre"
)

// SecureOpts defines options for the Secure middleware.
type SecureOpts struct {
	RedirectHTTP          bool
	XSSProtection         string
	FrameOptions          string
	ContentTypeOptions    string
	HSTSMaxAge            int
	HSTSIncludeSubdomains bool
	ContentSecurityPolicy string
	PublicKeyPins         string
}

var defaultSecureOpts = &SecureOpts{
	XSSProtection:      "1; mode=block",
	FrameOptions:       "SAMEORIGIN",
	ContentTypeOptions: "nosniff",
}

// Secure defines middleware for specifying secure headers.
func Secure(opts ...*SecureOpts) fibre.MiddlewareFunc {
	return func(h fibre.HandlerFunc) fibre.HandlerFunc {
		return func(c *fibre.Context) error {

			var config *SecureOpts

			switch len(opts) {
			case 0:
				config = defaultSecureOpts
			default:
				config = opts[0]
			}

			if config.RedirectHTTP && (!c.IsTLS() || c.Request().Header().Get(fibre.HeaderXForwardedProto) == "http") {
				h := c.Request().Host
				u := c.Request().RequestURI
				return c.Redirect(301, "https://"+h+u)
			}

			// This is a socket
			if c.IsSocket() {
				return h(c)
			}

			if config.XSSProtection != "" {
				c.Response().Header().Set(fibre.HeaderXXSSProtection, config.XSSProtection)
			}

			if config.FrameOptions != "" {
				c.Response().Header().Set(fibre.HeaderXFrameOptions, config.FrameOptions)
			}

			if config.ContentTypeOptions != "" {
				c.Response().Header().Set(fibre.HeaderXContentTypeOptions, config.ContentTypeOptions)
			}

			if config.ContentSecurityPolicy != "" {
				c.Response().Header().Set(fibre.HeaderContentSecurityPolicy, config.ContentSecurityPolicy)
			}

			if config.PublicKeyPins != "" {
				c.Response().Header().Set(fibre.HeaderPublicKeyPins, config.PublicKeyPins)
			}

			if (c.IsTLS() || c.Request().Header().Get(fibre.HeaderXForwardedProto) == "https") && config.HSTSMaxAge != 0 {
				if config.HSTSIncludeSubdomains {
					c.Response().Header().Set(fibre.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d; includeSubdomains", config.HSTSMaxAge))
				} else {
					c.Response().Header().Set(fibre.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d", config.HSTSMaxAge))
				}
			}

			return h(c)

		}
	}
}
