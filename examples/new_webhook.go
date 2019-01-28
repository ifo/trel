package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ifo/trel"
)

func main() {
	username := flag.String("u", "username", "your trello username or id")
	apiKey := flag.String("k", "apikey", "your trello api key")
	token := flag.String("t", "token", "your token - https://trello.com/app-key")
	listID := flag.String("l", "list id", "the list you want to watch")
	callbackURL := flag.String("cb", "http://example.com", "your callback url")
	flag.Parse()

	client := trel.New(
		nil,       // Default http client
		*username, // Your username
		*apiKey,   // Your api key
		*token,    // Your token (https://trello.com/app-key)
	)

	webhook, err := client.NewWebhook("webhook description: list watcher", *callbackURL, *listID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(webhook.Description)
}
