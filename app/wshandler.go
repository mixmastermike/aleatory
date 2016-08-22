package main

import (
	"net/http"

	"github.com/mixmastermike/aleatory/client"
	"github.com/mixmastermike/aleatory/provider"

	"github.com/gorilla/websocket"
)

var (
	upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

// WebSocket handler struct
type wsHandler struct {
	factory provider.Provider
}

// Handle an HTTP request
func (wsh wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Upgrade the HTTP request to a WebSocket
	sock, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Create a new client connection
	c := client.NewConnection(sock)
	// .. and register it into the data provider
	wsh.factory.Register(c)
	// Close the client connection when the http request is closed
	defer func() { wsh.factory.Unregister() }()
	// Begin procesing the communication streams
	go c.Writer()
	c.Reader()
}
