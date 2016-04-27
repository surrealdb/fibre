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

func NewURL(tls bool, uri string) *URL {

	var user string
	var pass string
	var host string

	var part *url.URL

	if tls == true {
		part, _ = url.Parse("https://" + uri)
	}

	if tls == false {
		part, _ = url.Parse("http://" + uri)
	}

	if part.User != nil {
		user = part.User.Username()
		pass, _ = part.User.Password()
	}

	host = strings.SplitN(part.Host, ":", 2)[0]

	return &URL{
		Scheme:   part.Scheme,
		User:     user,
		Pass:     pass,
		Host:     host,
		Path:     part.Path,
		Query:    part.RawQuery,
		Fragment: part.Fragment,
	}

}
