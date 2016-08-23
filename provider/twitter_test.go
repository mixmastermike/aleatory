package provider

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/mixmastermike/aleatory/client"
	"github.com/stretchr/testify/assert"
)

type cstHandler struct {
	t *testing.T
	C *client.Connection
}

type cstServer struct {
	*httptest.Server
	URL string
}

func (h cstHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.t.Logf("Upgrade: %v", err)
		return
	}
	defer ws.Close()

	op, rd, err := ws.NextReader()
	if err != nil {
		h.t.Logf("NextReader: %v", err)
		return
	}
	wr, err := ws.NextWriter(op)
	if err != nil {
		h.t.Logf("NextWriter: %v", err)
		return
	}
	// Copy incoming messages back to outgoing
	if _, err = io.Copy(wr, rd); err != nil {
		h.t.Logf("NextWriter: %v", err)
		return
	}
	if err := wr.Close(); err != nil {
		h.t.Logf("Close: %v", err)
		return
	}

}

func newWsServer(t *testing.T) *cstServer {
	var s cstServer
	s.Server = httptest.NewServer(cstHandler{t: t})
	s.URL = "ws" + strings.TrimPrefix(s.Server.URL, "http")
	return &s
}

func TestTwitterBroadcast(t *testing.T) {
	s := newWsServer(t)
	defer s.Close()

	dialer := websocket.Dialer{ReadBufferSize: 1024, WriteBufferSize: 1024}
	ws, _, err := dialer.Dial(s.URL, nil)
	assert.Nil(t, err)
	assert.NotNil(t, ws)
	defer ws.Close()

	sock := client.NewConnection(ws)
	p := &TwitterProvider{
		c: sock,
	}
	go sock.Writer()
	go sock.Reader()

	p.Broadcast([]byte("the test string"))
	// Read the message back off the socket
	_, ret, err := ws.ReadMessage()
	assert.Nil(t, err)
	assert.Equal(t, []byte("{\"type\":\"0\",\"weight\":15}"), ret)
}

func TestTwitterGenerateBroadcast(t *testing.T) {
	p := &TwitterProvider{}

	ret := p.generateBroadcast([]byte(""))
	assert.Equal(t, ret, []byte("{\"type\":\"0\",\"weight\":0}"))
	ret = p.generateBroadcast([]byte("the test tweet"))
	assert.Equal(t, ret, []byte("{\"type\":\"0\",\"weight\":14}"))
	ret = p.generateBroadcast([]byte("RT @me the test tweet"))
	assert.Equal(t, ret, []byte("{\"type\":\"1\",\"weight\":21}"))
}

func TestTwitterIsRetweet(t *testing.T) {
	p := &TwitterProvider{}

	assert.False(t, p.isReTweet([]byte{}))
	assert.False(t, p.isReTweet([]byte("1234")))
	assert.True(t, p.isReTweet([]byte("RT @")))
	assert.False(t, p.isReTweet([]byte("RT 1")))
	assert.False(t, p.isReTweet([]byte("some tweet about something")))
	assert.True(t, p.isReTweet([]byte("RT @somename some tweet about something")))
}

func TestTwitterStop(t *testing.T) {
	p := NewTwitterProvider(&TwitterConfig{})

	go func() { p.Stop() }()
	ret := <-p.done
	assert.True(t, ret)
}
