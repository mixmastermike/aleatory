package provider

import (
	"testing"

	"github.com/mixmastermike/aleatory/client"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestGenerateBroadcast(t *testing.T) {
	c := client.NewConnection(&websocket.Conn{})
	p := &TwitterProvider{
		c: c,
	}

	ret := []byte("foo")
	msg := []byte("theteststring")

	p.Broadcast(msg)
	ret = <-c.send
	assert.Equal(t, ret, msg)
}

func TestIsRetweet(t *testing.T) {
	p := &TwitterProvider{}
	assert.False(t, p.isReTweet([]byte{}))
	assert.False(t, p.isReTweet([]byte("1234")))
	assert.True(t, p.isReTweet([]byte("RT @")))
	assert.False(t, p.isReTweet([]byte("RT 1")))
	assert.False(t, p.isReTweet([]byte("some tweet about something")))
	assert.True(t, p.isReTweet([]byte("RT @somename some tweet about something")))
}
