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

// +build !appengine

package fibre

import (
	"log"
	"net/http"

	"gopkg.in/tylerb/graceful.v1"
)

// Run runs the server and handles http requests.
func (f *Fibre) Run(opts HTTPOptions, files ...string) {

	var err error
	var s *graceful.Server

	w := f.logger.Writer()
	defer w.Close()

	switch v := opts.(type) {
	case string:
		s = &graceful.Server{
			Timeout: f.wait,
			Server: &http.Server{
				Addr:         v,
				Handler:      f,
				ReadTimeout:  f.rtimeout,
				WriteTimeout: f.wtimeout,
				ErrorLog:     log.New(w, "", 0),
			},
		}
	case *http.Server:
		s = &graceful.Server{
			Timeout: f.wait,
			Server:  v,
		}
		s.Server.Handler = f
	case *graceful.Server:
		s = v
		s.Server.Handler = f
	}

	if len(files) != 2 {
		err = s.ListenAndServe()
	}

	if len(files) == 2 {
		err = s.ListenAndServeTLS(files[0], files[1])
	}

	if err != nil {
		f.Logger().Fatal(err)
	}

}
