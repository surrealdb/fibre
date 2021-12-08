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

package fibre

import (
	"encoding/xml"
	"github.com/gorilla/websocket"
	"github.com/ugorji/go/codec"
	"net/http"
)

// Socket wraps an websocket.Conn
type Client struct {
	*websocket.Conn
}

// NewClient creates a new instance of Response.
func NewClient(url string, protocols []string) (*Client, error) {

	dialer := &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		Subprotocols:      protocols,
		EnableCompression: true,
	}

	con, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return &Client{con}, nil

}

func (c *Client) Close() error {
	return c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

// Read reads a message from the socket.
func (c *Client) Read() (int, []byte, error) {
	return c.Conn.ReadMessage()
}

// ReadXML reads a xml message from the socket.
func (c *Client) ReadXML(v interface{}) (err error) {
	_, r, err := c.NextReader()
	if err != nil {
		return err
	}
	return xml.NewDecoder(r).Decode(v)
}

// ReadJSON reads a json message from the socket.
func (c *Client) ReadJSON(v interface{}) (err error) {
	_, r, err := c.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &jh).Decode(v)
}

// ReadCBOR reads a cbor message from the socket.
func (c *Client) ReadCBOR(v interface{}) (err error) {
	_, r, err := c.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &ch).Decode(v)
}

// ReadPACK reads a msgpack message from the socket.
func (c *Client) ReadPACK(v interface{}) (err error) {
	_, r, err := c.NextReader()
	if err != nil {
		return err
	}
	return codec.NewDecoder(r, &mh).Decode(v)
}

// Send sends a response to the socket.
func (c *Client) Send(t int, data []byte) (err error) {
	return c.Conn.WriteMessage(t, data)
}

// SendText sends a text response with status code.
func (c *Client) SendText(data string) (err error) {
	return c.Conn.WriteMessage(websocket.TextMessage, []byte(data))
}

// SendXML sends a xml response with status code.
func (c *Client) SendXML(data interface{}) (err error) {
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if data != nil {
		xml.NewEncoder(w).Encode(data)
	}
	return w.Close()
}

// SendJSON sends a json response with status code.
func (c *Client) SendJSON(data interface{}) (err error) {
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &jh).Encode(data)
	}
	return w.Close()
}

// SendCBOR sends a cbor response with status code.
func (c *Client) SendCBOR(data interface{}) (err error) {
	w, err := c.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &ch).Encode(data)
	}
	return w.Close()
}

// SendPACK sends a msgpack response with status code.
func (c *Client) SendPACK(data interface{}) (err error) {
	w, err := c.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	if data != nil {
		codec.NewEncoder(w, &mh).Encode(data)
	}
	return w.Close()
}

func (c *Client) Rpc() (chan<- *RPCRequest, <-chan *RPCResponse, chan error) {

	send := make(chan *RPCRequest)
	recv := make(chan *RPCResponse)
	quit := make(chan error, 1)
	exit := make(chan int, 1)
	kind := c.Subprotocol()

	go func() {
	loop:
		for {
			select {
			case <-exit:
				break loop
			default:

				var err error
				var req RPCResponse

				switch kind {
				case "json":
					err = c.ReadJSON(&req)
				case "cbor":
					err = c.ReadCBOR(&req)
				case "pack":
					err = c.ReadPACK(&req)
				}

				if err != nil {
					c.Close()
					quit <- err
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
			case res := <-send:

				var err error

				switch kind {
				case "json":
					err = c.SendJSON(res)
				case "cbor":
					err = c.SendCBOR(res)
				case "pack":
					err = c.SendPACK(res)
				}

				if err != nil {
					c.Close()
					quit <- err
					exit <- 0
					break loop
				}

			}
		}
	}()

	return send, recv, quit

}
