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
	listToID := flag.String("l", "to list id", "the list id to move your card to")
	cardID := flag.String("c", "card id", "the card to move")
	flag.Parse()

	client := trel.New(
		nil,     // Default http client
		*apiKey, // Your api key
		*token,  // Your token (https://trello.com/app-key)
	)

	card, err := client.Card(*cardID)
	if err != nil {
		log.Fatal(err)
	}

	// Move the card to the new list
	err = card.Move(*listToID)
	if err != nil {
		log.Fatal(err)
	}

	// Should print the listToID
	fmt.Println(card.IDList)
}
