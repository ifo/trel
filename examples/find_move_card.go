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
	listFromID := flag.String("lf", "from list id", "the list the card is on")
	listToID := flag.String("l", "to list id", "the list id to move your card to")
	cardName := flag.String("name", "card name", "the name of the card to move")
	flag.Parse()

	client := trel.New(
		nil,     // Default http client
		*apiKey, // Your api key
		*token,  // Your token (https://trello.com/app-key)
	)

	fromList, err := client.List(*listFromID)
	if err != nil {
		log.Fatal(err)
	}

	card, err := fromList.FindCard(*cardName)
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
