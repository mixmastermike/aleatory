package client

import "github.com/gorilla/websocket"

// Connection is a structure holding a client connection
type Connection struct {
	// The websocket connection
	ws *websocket.Conn
	// Channel for outbound messages
	send chan []byte
	// Func for inbound messages to be sent to
	rec func([]byte)
}

// NewConnection is a helper function to create a new Connection object
func NewConnection(ws *websocket.Conn) *Connection {
	return &Connection{
		ws:   ws,
		send: make(chan []byte, 1024),
		rec:  nil,
	}
}

// Broadcast will send a message to the client
func (c *Connection) Broadcast(msg []byte) {
	c.send <- msg
}

// Reader processes incoming messages from the client
func (c *Connection) Reader() {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		if c.rec != nil {
			c.rec(message)
		}
	}
	c.Close()
}

// Writer processes outbound messages to the client
func (c *Connection) Writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.Close()
}

// Receive will attach a new handler to the receieve function
func (c *Connection) Receive(fnc func([]byte)) {
	c.rec = fnc
}

// Close will close the connection
func (c *Connection) Close() error {
	return c.ws.Close()
}
