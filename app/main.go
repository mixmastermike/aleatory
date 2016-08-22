package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mixmastermike/aleatory/provider"
)

var (
	flags          = flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey    = flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret = flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken    = flags.String("access-token", "", "Twitter Access Token")
	accessSecret   = flags.String("access-secret", "", "Twitter Access Secret")
	// Allow changing the listening address via flag
	addr = flag.String("addr", ":8080", "http service address")
)

func main() {
	flags.Parse(os.Args[1:])
	// Make sure we have everything we need
	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	// Handle page requests
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Render the page, passing in the value of this server for the WebSocket conection
		renderTemplate(w, "static/index.html", "index", map[string]string{"Host": req.Host})
	})
	// Handle the WebSocket requets
	http.Handle("/ws", wsHandler{
		factory: provider.NewTwitterProvider(&provider.TwitterConfig{
			ConsumerKey:    *consumerKey,
			ConsumerSecret: *consumerSecret,
			AccessToken:    *accessToken,
			AccessSecret:   *accessSecret,
		}),
	})
	// Handle the rest
	http.Handle("/js/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	http.Handle("/css/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	http.Handle("/img/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))
	http.Handle("/sounds/", http.StripPrefix("/", http.FileServer(http.Dir("static"))))

	fmt.Printf("Service listening on port %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
