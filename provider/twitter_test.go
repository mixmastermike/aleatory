package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
