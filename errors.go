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
)

const (
	rpcError        = "Unknown error"
	rpcParseError   = "Parse error"
	rpcInvalidError = "Invalid Request"
	rpcMethodError  = "Method not found"
	rpcParamsError  = "Invalid params"
)

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	code    int
	message string
	fields  map[string]interface{}
}

// NewHTTPError creates a new instance of HTTPError.
func NewHTTPError(code int, message ...string) (err *HTTPError) {

	err = &HTTPError{code: code}

	switch len(message) {
	case 1:
		err.message = message[0]
	case 0:
		err.message = http.StatusText(code)
	}

	return err

}

// Code returns code.
func (e *HTTPError) Code() int {
	return e.code
}

// Error returns message.
func (e *HTTPError) Error() string {
	return e.message
}

// Fields returns error fields.
func (e *HTTPError) Fields() map[string]interface{} {
	return e.fields
}

func (e *HTTPError) WithMessage(m string) *HTTPError {
	e.message = m
	return e
}

func (e *HTTPError) WithField(k string, v interface{}) *HTTPError {
	if e.fields == nil {
		e.fields = make(map[string]interface{})
	}
	e.fields[k] = v
	return e
}

func (e *HTTPError) WithFields(f map[string]interface{}) *HTTPError {
	if e.fields == nil {
		e.fields = make(map[string]interface{})
	}
	for k, v := range f {
		e.fields[k] = v
	}
	return e
}

// DefaultHTTPErrorHandler invokes the default HTTP error handler.
func (f *Fibre) defaultErrorHandler(err error, c *Context) {

	code := http.StatusInternalServerError
	mesg := http.StatusText(code)

	if he, ok := err.(*HTTPError); ok {
		code = he.Code()
		mesg = he.Error()
	}

	if !c.Response().Done() {
		http.Error(c.response, mesg, code)
	}

	f.Logger().Debugf("%v", err)

}
