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

package fibre

import (
	"net/url"
	"strings"
)

// URL defines a parsed url
type URL struct {
	Scheme   string
	User     string
	Pass     string
	Host     string
	Path     string
	Query    string
	Fragment string
}

func NewURL(uri string) *URL {

	var user string
	var pass string
	var host string

	parsed, _ := url.Parse(uri)

	if parsed.User != nil {
		user = parsed.User.Username()
		pass, _ = parsed.User.Password()
	}

	host = strings.SplitN(parsed.Host, ":", 2)[0]

	return &URL{
		Scheme:   parsed.Scheme,
		User:     user,
		Pass:     pass,
		Host:     host,
		Path:     parsed.Path,
		Query:    parsed.RawQuery,
		Fragment: parsed.Fragment,
	}

}
