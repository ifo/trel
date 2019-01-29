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
	listID := flag.String("l", "list id", "your list's id")
	flag.Parse()

	client := trel.New(
		nil,     // Default http client
		*apiKey, // Your api key
		*token,  // Your token (https://trello.com/app-key)
	)

	list, err := client.List(*listID)
	if err != nil {
		log.Fatal(err)
	}

	cards, err := list.Cards()
	if err != nil {
		log.Fatal(err)
	}
	for _, card := range cards {
		fmt.Println(card.Name)
	}
}
