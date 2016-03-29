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
	"net/http"
	"time"
)

// Request wraps an http.Request
type Request struct {
	*http.Request
	fibre *Fibre
	start time.Time
}

// NewRequest creates a new instance of Response.
func NewRequest(i *http.Request, f *Fibre) *Request {
	return &Request{i, f, time.Now()}
}

// Size returns the current size, in bytes, of the request.
func (r *Request) Size() int64 {
	return r.Request.ContentLength
}

// Start returns the current size, in bytes, of the request.
func (r *Request) Start() time.Time {
	return r.start
}

// Reader returns the http.Request instance for this Request.
func (r *Request) Reader() *http.Request {
	return r.Request
}

// SetReader sets the http.Request instance for this Request.
func (r *Request) SetReader(w *http.Request) {
	r.Request = w
}

// Header returns the header map values for this Request,
func (r *Request) Header() http.Header {
	return r.Request.Header
}

func (r *Request) reset(i *http.Request, f *Fibre) {
	r.fibre = f
	r.Request = i
	r.start = time.Now()
}
