package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateBroadcast(t *testing.T) {
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
