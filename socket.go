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
	"time"

	"encoding/xml"
	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
)

var ping = []byte(nil)

// Socket wraps an websocket.Conn
type Socket struct {
	*websocket.Conn
	context *Context
	fibre   *Fibre
	notify  chan<- *RPCNotification
}

// NewSocket creates a new instance of Response.
func NewSocket(i *websocket.Conn, c *Context, f *Fibre) *Socket {
	return &Socket{i, c, f, nil}
}

func (s *Socket) err(err error) error {
	if websocket.IsCloseError(err, 1000, 1001, 1005) {
		return nil
	}
	return NewHTTPError(400)
}

func (s *Socket) rpc() (chan<- *RPCResponse, <-chan *RPCRequest, chan error) {

	s.SetPongHandler(func(msg string) error {
		return nil
	})

	noti := make(chan *RPCNotification)
	send := make(chan *RPCResponse)
	recv := make(chan *RPCRequest)
	quit := make(chan error, 1)
	exit := make(chan int, 1)
	kind := s.Subprotocol()

	s.notify = noti

	go func() {
	loop:
		for {
			select {
			case <-exit:
				break loop
			case <-time.After(25 * time.Second):

				var err error

				var dead = time.Now().Add(5 * time.Second)

				err = s.WriteControl(websocket.PingMessage, ping, dead)

				if err != nil {
					s.Close(websocket.CloseNoStatusReceived)
					quit <- s.err(err)
					exit <- 0
					break loop
				}

			}
		}
	}()

	go func() {
	loop:
		for {
			select {
			case <-exit:
				break loop
			default:

				var err error
				var req RPCRequest

				switch kind {
				case "json":
					err = s.ReadJSON(&req)
				case "cbor":
					err = s.ReadCBOR(&req)
				case "pack":
					err = s.ReadPACK(&req)
				}

				if err != nil {
					s.Close(websocket.CloseUnsupportedData)
					quit <- s.err(err)
					exit <- 0
					break loop
				}

				recv <- &req

			}
		}
	}()

	go func() {
	loop:
		for {
			select {
			case <-exit:
				break loop
			case res := <-noti:

				var err error

				switch kind {
				case "json":
					err = s.SendJSON(res)
				case "cbor":
					err = s.SendCBOR(res)
				case "pack":
					err = s.SendPACK(res)
				}

				if err != nil {
					s.Close(websocket.CloseUnsupportedData)
					quit <- s.err(err)
					exit <- 0
					break loop
				}

			case res := <-send:

				var err error

				switch kind {
				case "json":
					err = s.SendJSON(res)
				case "cbor":
					err = s.SendCBOR(res)
				case "pack":
					err = s.SendPACK(res)
				}

				if err != nil {
					s.Close(websocket.CloseUnsupportedData)
					quit <- s.err(err)
					exit <- 0
					break loop
				}

			}
		}
	}()

	return send, recv, quit

}

func (s *Socket) Close(code int) error {
	return s.Conn.Close()
}

func (s *Socket) Notify(val *RPCNotification) {
	if s.notify != nil {
		s.notify <- val
	}
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

// ReadPACK reads a msgpack message from the socket.
func (s *Socket) ReadPACK(v interface{}) (err error) {
	_, r, err := s.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &mh).Decode(v)
}

// Send sends a response to the socket.
func (s *Socket) Send(t int, data []byte) (err error) {
	return s.Conn.WriteMessage(t, data)
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
		if err = xml.NewEncoder(w).Encode(data); err != nil {
			return err
		}
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
		if err = codec.NewEncoder(w, &jh).Encode(data); err != nil {
			return err
		}
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
		if err = codec.NewEncoder(w, &ch).Encode(data); err != nil {
			return err
		}
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
		if err = codec.NewEncoder(w, &mh).Encode(data); err != nil {
			return err
		}
	}
	return w.Close()
}
