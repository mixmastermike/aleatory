package provider

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mixmastermike/aleatory/client"
)

const (
	messageTypeTweet   string = "0"
	messageTypeRetweet string = "1"
	ucR                byte   = byte('R')
	ucT                byte   = byte('T')
	space              byte   = byte(' ')
	at                 byte   = byte('@')
)

// TwitterProvider is a data provider using Twitter stream data
type TwitterProvider struct {
	Provider

	// The Twitter client
	client *twitter.Client
	// The client WebSocket connection
	c *client.Connection
	// A Channel to inform to generate() to stop
	done chan bool
}

// TwitterConfig is the configuration required to access the Twitter API
type TwitterConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

// TwitterGenerateParams is the set of parameters that inform data generation
type TwitterGenerateParams struct {
	GenerateParams

	search string
}

var errClientAlreadyRegistered = errors.New("Client already registered")

// NewTwitterProvider will return a new TwitterProvider based on the given
// configuration
func NewTwitterProvider(conf *TwitterConfig) *TwitterProvider {

	config := oauth1.NewConfig(conf.ConsumerKey, conf.ConsumerSecret)
	token := oauth1.NewToken(conf.AccessToken, conf.AccessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	return &TwitterProvider{
		client: client,
		done:   make(chan bool),
	}
}

// Broadcast a message to the connected clients
func (t *TwitterProvider) Broadcast(msg []byte) {
	if t.c != nil {
		t.c.Broadcast(t.generateBroadcast(msg))
	}
}

// Register a client with this data provider. Register will return an error if
// a client is already registered.  You must Unregister existing clients first.
func (t *TwitterProvider) Register(c *client.Connection) error {
	if t.c != nil {
		return errClientAlreadyRegistered
	}
	t.c = c
	// Attach the listener
	t.c.Receive(t.listenHandler)
	return nil
}

// Unregister a client with this data provider
func (t *TwitterProvider) Unregister() (*client.Connection, error) {

	c := t.c
	if c != nil {
		// Inform the client to close
		err := t.c.Close()
		if err != nil {
			return nil, err
		}
		t.c = nil
		// No need to keep on producing ..
		t.done <- true
	}
	return c, nil
}

// Generate starts the provider producing data
func (t *TwitterProvider) Generate(p GenerateParams) {

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	// Handle tweets
	demux.Tweet = func(tweet *twitter.Tweet) {
		t.Broadcast([]byte(tweet.Text))
	}
	// Set up the Twitter stream
	stream := t.streamSetup(p)
	// Reset the done chan
	t.done = make(chan bool)

MessageLoop:
	for {
		// process messages until we receieve a signal
		select {
		case <-t.done:
			break MessageLoop
		case message := <-stream.Messages:
			demux.Handle(message)
		}
	}
	stream.Stop()
}

// Stop the provider from generating data
func (t *TwitterProvider) Stop() {
	t.done <- true
}

// Set up the Twitter data stream
func (t *TwitterProvider) streamSetup(p GenerateParams) *twitter.Stream {

	// Set up the stream criteria
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{p.(*TwitterGenerateParams).search},
		StallWarnings: twitter.Bool(true),
	}

	// Set up the stream
	stream, err := t.client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	return stream
}

// Generate the broadcast message to be send to the WebSocket client
func (t *TwitterProvider) generateBroadcast(m []byte) []byte {

	mtyp := messageTypeTweet
	if t.isReTweet(m) {
		mtyp = messageTypeRetweet
	}

	// Format a message object then send back as JSON bytes
	msg := client.Message{
		Type:   mtyp,
		Weight: len(m),
	}
	b, _ := json.Marshal(msg)
	return b
}

func (t *TwitterProvider) isReTweet(m []byte) bool {
	if len(m) < 4 {
		return false
	}
	// We want a prefix of "RT @" to denote a retweet
	return m[0] == ucR &&
		m[1] == ucT &&
		m[2] == space &&
		m[3] == at
}

// Handle messages from the WebSocket client
func (t *TwitterProvider) listenHandler(m []byte) {
	msg := client.Command{}
	json.Unmarshal(m, &msg)

	switch strings.ToLower(msg.Command) {
	case "search":
		// Start generating some data
		go t.Generate(&TwitterGenerateParams{
			search: msg.Value.(string),
		})
	default:
		// noop
	}
}
