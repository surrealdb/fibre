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

// +build !appengine

package fibre

import (
	"log"
	"net/http"

	"github.com/ory/graceful"
)

// Run runs the server and handles http requests.
func (f *Fibre) Run(a string, files ...string) {

	var err error

	w := f.logger.Writer()
	defer w.Close()

	s := graceful.WithDefaults(&http.Server{
		Addr:         a,
		Handler:      f,
		IdleTimeout:  f.itimeout,
		ReadTimeout:  f.rtimeout,
		WriteTimeout: f.wtimeout,
		ErrorLog:     log.New(w, "", 0),
	})

	if len(files) != 2 {
		err = graceful.Graceful(func() error {
			return s.ListenAndServe()
		}, s.Shutdown)
	}

	if len(files) == 2 {
		err = graceful.Graceful(func() error {
			return s.ListenAndServeTLS(files[0], files[1])
		}, s.Shutdown)
	}

	if err != nil {
		f.Logger().Fatal(err)
	}

}
