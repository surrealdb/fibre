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
	"encoding/xml"
	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
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

func (s *Socket) err(err error) error {
	if websocket.IsUnexpectedCloseError(err, 1000, 1001, 1005) {
		return err
	}
	return nil
}

func (s *Socket) rpc(format string) (chan<- *RPCResponse, <-chan *RPCRequest, chan error) {

	send := make(chan *RPCResponse)
	recv := make(chan *RPCRequest)
	quit := make(chan error, 1)
	exit := make(chan int, 1)

	go func() {
	loop:
		for {
			select {
			case <-exit:
				break loop
			default:

				var e error
				var v RPCRequest

				switch format {
				default:
					e = s.ReadJSON(&v)
				case "JSON":
					e = s.ReadJSON(&v)
				case "CBOR":
					e = s.ReadCBOR(&v)
				case "BINC":
					e = s.ReadBINC(&v)
				case "PACK":
					e = s.ReadPACK(&v)
				}

				if e != nil {
					s.Close()
					quit <- s.err(e)
					exit <- 0
					break loop
				}

				recv <- &v

			}
		}
	}()

	go func() {
	loop:
		for {
			select {
			case <-exit:
				break loop
			case v := <-send:

				var e error

				switch format {
				default:
					e = s.SendJSON(v)
				case "JSON":
					e = s.SendJSON(v)
				case "CBOR":
					e = s.SendCBOR(v)
				case "BINC":
					e = s.SendBINC(v)
				case "PACK":
					e = s.SendPACK(v)
				}

				if e != nil {
					s.Close()
					quit <- s.err(e)
					exit <- 0
					break loop
				}

			}
		}
	}()

	return send, recv, quit

}

// Read reads a message from the socket.
func (s *Socket) Read() (int, []byte, error) {
	return s.Conn.ReadMessage()
}

// ReadXML reads a xml message from the socket.
func (s *Socket) ReadXML(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return xml.NewDecoder(r).Decode(v)
}

// ReadJSON reads a json message from the socket.
func (s *Socket) ReadJSON(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &jh).Decode(v)
}

// ReadCBOR reads a cbor message from the socket.
func (s *Socket) ReadCBOR(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &ch).Decode(v)
}

// ReadBINC reads a binc message from the socket.
func (s *Socket) ReadBINC(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &bh).Decode(v)
}

// ReadPACK reads a msgpack message from the socket.
func (s *Socket) ReadPACK(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &mh).Decode(v)
}

// SendText sends a text response with status code.
func (s *Socket) SendText(data string) (err error) {
	return s.Conn.WriteMessage(websocket.TextMessage, []byte(data))
}

// SendXML sends a xml response with status code.
func (s *Socket) SendXML(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if data != nil {
		xml.NewEncoder(w).Encode(data)
	}
	return w.Close()
}

// SendJSON sends a json response with status code.
func (s *Socket) SendJSON(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &jh).Encode(data)
	}
	return w.Close()
}

// SendCBOR sends a cbor response with status code.
func (s *Socket) SendCBOR(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &ch).Encode(data)
	}
	return w.Close()
}

// SendBINC sends a binc response with status code.
func (s *Socket) SendBINC(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &bh).Encode(data)
	}
	return w.Close()
}

// SendPACK sends a msgpack response with status code.
func (s *Socket) SendPACK(data interface{}) (err error) {
	w, err := s.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &mh).Encode(data)
	}
	return w.Close()
}
