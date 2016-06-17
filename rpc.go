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

	"github.com/gorilla/websocket"
)

// RPCError represents a jsonrpc error
type RPCError struct {
	Code    int    `json:"code" msgpack:"code"`
	Message string `json:"message" msgpack:"message"`
}

// RPCRequest represents an incoming jsonrpc request
type RPCRequest struct {
	ID     string        `json:"id" msgpack:"id"`
	Method string        `json:"method" msgpack:"method"`
	Params []interface{} `json:"params" msgpack:"params"`
}

// RPCResponse represents an outgoing jsonrpc response
type RPCResponse struct {
	ID     string      `json:"id" msgpack:"id"`
	Error  *RPCError   `json:"error,omitempty" msgpack:"error,omitempty"`
	Result interface{} `json:"result,omitempty" msgpack:"result,omitempty"`
}

// Rpc adds a route > handler to the router for a jsonrpc endpoint.
func (f *Fibre) Rpc(p string, i interface{}) {

	f.router.Add(POST, p, func(c *Context) (err error) {
		req := &RPCRequest{}
		c.Bind(req)
		res := rpc(req, c, i)
		return c.Send(200, res)
	})

	f.router.Add(GET, p, func(c *Context) (err error) {

		if err = c.Upgrade(); err != nil {
			return
		}

		for {
			req := &RPCRequest{}
			if err = c.Socket().ReadJSON(req); err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					break
				}
			}
			res := rpc(req, c, i)
			if err = c.Socket().SendJSON(res); err != nil {
				if _, ok := err.(*websocket.CloseError); ok {
					break
				}
			}
		}

		return nil

	})

}

func rpc(req *RPCRequest, c *Context, i interface{}) interface{} {

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

	if req.ID == "" {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32600,
				Message: "Invalid Request",
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

	_, ok := ins.Type().MethodByName(req.Method)
	if !ok {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32601,
				Message: "Method not found",
			},
		}
	}

	fnc := ins.MethodByName(req.Method)

	cnti := fnc.Type().NumIn()
	if cnti != len(req.Params)+1 {
		return &RPCResponse{
			ID: req.ID,
			Error: &RPCError{
				Code:    -32602,
				Message: "Invalid params",
			},
		}
	}

	cnto := fnc.Type().NumOut()
	if cnto != 2 {
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

	for k, v := range req.Params {
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

func arg(fnc reflect.Value, k int, v interface{}) (reflect.Value, error) {

	typf := fnc.Type().In(k + 1)
	str := v.(string)

	switch typf.Kind() {
	default:
		return reflect.ValueOf(str), nil

	case reflect.String:
		return reflect.ValueOf(str), nil

	case reflect.Bool:
		cnv, err := strconv.ParseBool(str)
		return reflect.ValueOf(cnv), err

	case reflect.Float32:
		cnv, err := strconv.ParseFloat(str, 32)
		return reflect.ValueOf(float32(cnv)), err
	case reflect.Float64:
		cnv, err := strconv.ParseFloat(str, 64)
		return reflect.ValueOf(float64(cnv)), err

	case reflect.Int:
		cnv, err := strconv.ParseInt(str, 10, 0)
		return reflect.ValueOf(int(cnv)), err
	case reflect.Int8:
		cnv, err := strconv.ParseInt(str, 10, 8)
		return reflect.ValueOf(int8(cnv)), err
	case reflect.Int16:
		cnv, err := strconv.ParseInt(str, 10, 16)
		return reflect.ValueOf(int16(cnv)), err
	case reflect.Int32:
		cnv, err := strconv.ParseInt(str, 10, 32)
		return reflect.ValueOf(int32(cnv)), err
	case reflect.Int64:
		cnv, err := strconv.ParseInt(str, 10, 64)
		return reflect.ValueOf(int64(cnv)), err

	case reflect.Uint:
		cnv, err := strconv.ParseUint(str, 10, 0)
		return reflect.ValueOf(uint(cnv)), err
	case reflect.Uint8:
		cnv, err := strconv.ParseUint(str, 10, 8)
		return reflect.ValueOf(uint8(cnv)), err
	case reflect.Uint16:
		cnv, err := strconv.ParseUint(str, 10, 16)
		return reflect.ValueOf(uint16(cnv)), err
	case reflect.Uint32:
		cnv, err := strconv.ParseUint(str, 10, 32)
		return reflect.ValueOf(uint32(cnv)), err
	case reflect.Uint64:
		cnv, err := strconv.ParseUint(str, 10, 64)
		return reflect.ValueOf(uint64(cnv)), err

	}

}
