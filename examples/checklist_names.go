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
	cardID := flag.String("c", "card id", "your card with checklists's id")
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

	checklists, err := card.Checklists()
	if err != nil {
		log.Fatal(err)
	}
	for _, checklist := range checklists {
		fmt.Println(checklist.Name)
	}
}
