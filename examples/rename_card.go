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
		*username, // Your username
		*apiKey,   // Your api key
		*token,    // Your token (https://trello.com/app-key)
	)

	list, err := client.List(*listID)
	if err != nil {
		log.Fatal(err)
	}

	// Add the card to the end of the list's cards
	newCard, err := list.NewCard("card name", "card description", "bottom")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(newCard.ID)

	err = newCard.Rename("new card name")
	if err != nil {
		log.Fatal(err)
	}

	// Should print "new card name - card description"
	fmt.Println(newCard.Name + " - " + newCard.Description)
}
