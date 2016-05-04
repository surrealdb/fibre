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
	"encoding/json"
	"encoding/xml"

	"github.com/gorilla/websocket"
	"gopkg.in/vmihailenco/msgpack.v2"
)

// Socket wraps an websocket.Conn
type Socket struct {
	*websocket.Conn
	context *Context
	fibre   *Fibre
}

// NewSocket creates a new instance of Response.
func NewSocket(i *websocket.Conn, c *Context, f *Fibre) *Socket {
	return &Socket{i, c, f}
}

// Read reads a message from the socket.
func (s *Socket) Read() (int, []byte, error) {
	return s.Conn.ReadMessage()
}

// Text sends a text response with status code.
func (s *Socket) Text(data string) (err error) {
	s.Conn.WriteMessage(websocket.TextMessage, []byte(data))
	return
}

// HTML sends an html response with status code.
func (s *Socket) HTML(data string) (err error) {
	s.Conn.WriteMessage(websocket.TextMessage, []byte(data))
	return
}

// XML sends a xml response with status code.
func (s *Socket) XML(data interface{}) (err error) {
	done, err := xml.Marshal(data)
	if err != nil {
		return err
	}
	s.Conn.WriteMessage(websocket.TextMessage, done)
	return
}

// JSON sends a json response with status code.
func (s *Socket) JSON(data interface{}) (err error) {
	done, err := json.Marshal(data)
	if err != nil {
		return err
	}
	s.Conn.WriteMessage(websocket.TextMessage, done)
	return
}

// PACK sends a msgpack response with status code.
func (s *Socket) PACK(data interface{}) (err error) {
	done, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	s.Conn.WriteMessage(websocket.BinaryMessage, done)
	return
}

// Send sends the relevant response depending on the request type.
func (s *Socket) Send(data interface{}) (err error) {
	switch s.context.Type() {
	default:
		return s.JSON(data)
	case "application/xml":
		return s.XML(data)
	case "application/json":
		return s.JSON(data)
	case "application/msgpack":
		return s.PACK(data)
	}
}
