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

func (s *Socket) ReadXML(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return xml.NewDecoder(r).Decode(v)
}

func (s *Socket) ReadJSON(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return json.NewDecoder(r).Decode(v)
}

func (s *Socket) ReadPACK(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return msgpack.NewDecoder(r).Decode(v)
}

// Text sends a text response with status code.
func (s *Socket) SendText(data string) (err error) {
	s.Conn.WriteMessage(websocket.TextMessage, []byte(data))
	return
}

// XML sends a xml response with status code.
func (s *Socket) SendXML(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.TextMessage)
	if data != nil {
		xml.NewEncoder(w).Encode(data)
	}
	return w.Close()
}

// JSON sends a json response with status code.
func (s *Socket) SendJSON(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.TextMessage)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
	return w.Close()
}

// PACK sends a msgpack response with status code.
func (s *Socket) SendPACK(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.BinaryMessage)
	if data != nil {
		msgpack.NewEncoder(w).Encode(data)
	}
	return w.Close()
}
