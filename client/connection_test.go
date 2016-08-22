package client

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	ws := &websocket.Conn{}
	c := NewConnection(ws)

	assert.NotNil(t, c)
	assert.Equal(t, c.ws, ws)
}

func TestConnectionBroadcast(t *testing.T) {
	ws := &websocket.Conn{}
	c := NewConnection(ws)

	ret := []byte("foo")
	msg := []byte("theteststring")
	c.send <- msg
	ret = <-c.send
	assert.Equal(t, ret, msg)
}
