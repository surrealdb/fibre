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
	"reflect"
	"strconv"
)

// RPCNull represents a null argument
type RPCNull struct{}

// RPCError represents a jsonrpc error
type RPCError struct {
	Code    int    `json:"code" msgpack:"code"`
	Message string `json:"message" msgpack:"message"`
}

// RPCRequest represents an incoming jsonrpc request
type RPCRequest struct {
	ID     interface{}   `json:"id" msgpack:"id"`
	Method string        `json:"method" msgpack:"method"`
	Params []interface{} `json:"params" msgpack:"params"`
}

// RPCResponse represents an outgoing jsonrpc response
type RPCResponse struct {
	ID     interface{} `json:"id" msgpack:"id"`
	Error  *RPCError   `json:"error,omitempty" msgpack:"error,omitempty"`
	Result interface{} `json:"result,omitempty" msgpack:"result,omitempty"`
}

// Rpc adds a route > handler to the router for a jsonrpc endpoint.
func (f *Fibre) Rpc(p string, i interface{}) {

	f.router.Add(POST, p, func(c *Context) (err error) {
		req := &RPCRequest{}
		c.Bind(req)
		if res := rpc(req, c, i); res != nil {
			return c.Send(200, res)
		}
		return c.Code(200)
	})

	f.router.Add(GET, p, func(c *Context) (err error) {

		if err = c.Upgrade(); err != nil {
			return
		}

		send, recv, quit := c.Socket().rpc(c.Query("format"))

		for {
			select {
			case err := <-quit:
				return err
			case req := <-recv:
				if res := rpc(req, c, i); res != nil {
					send <- res
				}
			}
		}

		return nil

	})

}

func rpc(req *RPCRequest, c *Context, i interface{}) (o *RPCResponse) {

	defer func() {

		if r := recover(); r != nil && req.ID != "" {

			o = &RPCResponse{
				ID: req.ID,
				Error: &RPCError{
					Code:    -32099,
					Message: "Unknown error",
				},
			}

			if err, ok := r.(error); ok {
				o.Error.Message = err.(error).Error()
			}

		}

		if req.ID == "" {
			o = nil
		}

	}()

	ins := reflect.ValueOf(i)

	if req == nil {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32700,
				Message: "Parse error",
			},
		}
	}

	if req.Method == "" {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32600,
				Message: "Invalid Request",
			},
		}
	}

	if _, ok := ins.Type().MethodByName(req.Method); !ok {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}

	fnc := ins.MethodByName(req.Method)

	if fnc.Type().NumIn()-1 < len(req.Params) {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	if fnc.Type().NumOut() != 2 {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	var args []reflect.Value

	args = append(args, reflect.ValueOf(c))

	for k := 0; k < fnc.Type().NumIn()-1; k++ {
		var v interface{}
		if k < len(req.Params) {
			v = req.Params[k]
		}
		val, err := arg(fnc, k, v)
		if err != nil {
			return &RPCResponse{
				ID: req.ID,
				Error: &RPCError{
					Code:    -32602,
					Message: "Invalid params",
				},
			}
		}
		args = append(args, val)
	}

	val := fnc.Call(args)
	res := val[0].Interface()
	err := val[1].Interface()

	if err == nil {
		return &RPCResponse{
			ID:     req.ID,
			Result: res,
		}
	}

	if err != nil {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32000,
				Message: err.(error).Error(),
			},
		}
	}

	return nil

}

func arg(fnc reflect.Value, k int, i interface{}) (reflect.Value, error) {

	a := fnc.Type().In(k + 1)

	switch a.Kind() {

	default:
		return reflect.ValueOf(i), nil

	case reflect.Interface:
		if i == nil {
			return reflect.ValueOf(new(RPCNull)), nil
		} else {
			return reflect.ValueOf(i), nil
		}

	case reflect.Map:
		if i == nil || a != reflect.TypeOf(i) {
			return reflect.MakeMap(a), nil
		} else {
			return reflect.ValueOf(i), nil
		}

	case reflect.Slice:
		if i == nil || a != reflect.TypeOf(i) {
			return reflect.MakeSlice(a, 0, 0), nil
		} else {
			return reflect.ValueOf(i), nil
		}

	case reflect.String:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(""), nil
		case bool:
			s := strconv.FormatBool(v)
			return reflect.ValueOf(s), nil
		case int64:
			s := strconv.FormatInt(v, 10)
			return reflect.ValueOf(s), nil
		case float64:
			s := strconv.FormatFloat(v, 'g', -1, 64)
			return reflect.ValueOf(s), nil
		case string:
			return reflect.ValueOf(v), nil
		}

	case reflect.Bool:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(false), nil
		case bool:
			return reflect.ValueOf(v), nil
		case string:
			c, e := strconv.ParseBool(v)
			return reflect.ValueOf(c), e
		}

	case reflect.Float32:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(float32(0)), nil
		case int64:
			return reflect.ValueOf(float32(v)), nil
		case float64:
			return reflect.ValueOf(float32(v)), nil
		case string:
			c, e := strconv.ParseFloat(v, 32)
			return reflect.ValueOf(c), e
		}

	case reflect.Float64:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(float64(0)), nil
		case int64:
			return reflect.ValueOf(float64(v)), nil
		case float64:
			return reflect.ValueOf(float64(v)), nil
		case string:
			c, e := strconv.ParseFloat(v, 64)
			return reflect.ValueOf(c), e
		}

	case reflect.Int:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(int(0)), nil
		case int64:
			return reflect.ValueOf(int(v)), nil
		case float64:
			return reflect.ValueOf(int(v)), nil
		case string:
			c, e := strconv.ParseInt(v, 10, 0)
			return reflect.ValueOf(c), e
		}

	case reflect.Int8:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(int8(0)), nil
		case int64:
			return reflect.ValueOf(int8(v)), nil
		case float64:
			return reflect.ValueOf(int8(v)), nil
		case string:
			c, e := strconv.ParseInt(v, 10, 8)
			return reflect.ValueOf(c), e
		}

	case reflect.Int16:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(int16(0)), nil
		case int64:
			return reflect.ValueOf(int16(v)), nil
		case float64:
			return reflect.ValueOf(int16(v)), nil
		case string:
			c, e := strconv.ParseInt(v, 10, 16)
			return reflect.ValueOf(c), e
		}

	case reflect.Int32:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(int32(0)), nil
		case int64:
			return reflect.ValueOf(int32(v)), nil
		case float64:
			return reflect.ValueOf(int32(v)), nil
		case string:
			c, e := strconv.ParseInt(v, 10, 32)
			return reflect.ValueOf(c), e
		}

	case reflect.Int64:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(int64(0)), nil
		case int64:
			return reflect.ValueOf(int64(v)), nil
		case float64:
			return reflect.ValueOf(int64(v)), nil
		case string:
			c, e := strconv.ParseInt(v, 10, 64)
			return reflect.ValueOf(c), e
		}

	case reflect.Uint:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(uint(0)), nil
		case int64:
			return reflect.ValueOf(uint(v)), nil
		case float64:
			return reflect.ValueOf(uint(v)), nil
		case string:
			c, e := strconv.ParseUint(v, 10, 0)
			return reflect.ValueOf(c), e
		}

	case reflect.Uint8:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(uint8(0)), nil
		case int64:
			return reflect.ValueOf(uint8(v)), nil
		case float64:
			return reflect.ValueOf(uint8(v)), nil
		case string:
			c, e := strconv.ParseUint(v, 10, 8)
			return reflect.ValueOf(c), e
		}

	case reflect.Uint16:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(uint16(0)), nil
		case int64:
			return reflect.ValueOf(uint16(v)), nil
		case float64:
			return reflect.ValueOf(uint16(v)), nil
		case string:
			c, e := strconv.ParseUint(v, 10, 16)
			return reflect.ValueOf(c), e
		}

	case reflect.Uint32:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(uint32(0)), nil
		case int64:
			return reflect.ValueOf(uint32(v)), nil
		case float64:
			return reflect.ValueOf(uint32(v)), nil
		case string:
			c, e := strconv.ParseUint(v, 10, 32)
			return reflect.ValueOf(c), e
		}

	case reflect.Uint64:
		switch v := i.(type) {
		default:
			return reflect.ValueOf(uint64(0)), nil
		case int64:
			return reflect.ValueOf(uint64(v)), nil
		case float64:
			return reflect.ValueOf(uint64(v)), nil
		case string:
			c, e := strconv.ParseUint(v, 10, 64)
			return reflect.ValueOf(c), e
		}

	}

}
