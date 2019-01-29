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
	flag.Parse()

	client := trel.New(
		nil,     // Default http client
		*apiKey, // Your api key
		*token,  // Your token (https://trello.com/app-key)
	)

	webhooks, err := client.Webhooks()
	if err != nil {
		log.Fatal(err)
	}
	if len(webhooks) == 0 {
		fmt.Println("No webhooks found")
		return
	}
	for _, webhook := range webhooks {
		fmt.Println(webhook.Description)
	}
}
