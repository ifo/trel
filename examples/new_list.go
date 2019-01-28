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
	boardID := flag.String("b", "board id", "your board's id")
	flag.Parse()

	client := trel.New(
		nil,       // Default http client
		*username, // Your username
		*apiKey,   // Your api key
		*token,    // Your token (https://trello.com/app-key)
	)

	board, err := client.Board(*boardID)
	if err != nil {
		log.Fatal(err)
	}

	// Add the list to the end of the board's lists
	newList, err := board.NewList("list name", "bottom")
	if err != nil {
		log.Fatal(err)
	}
	// Should print "list name"
	fmt.Println(newList.Name)
}
