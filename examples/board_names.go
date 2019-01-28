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
		nil,       // Default http client
		*username, // Your username
		*apiKey,   // Your api key
		*token,    // Your token (https://trello.com/app-key)
	)

	// Get the first board
	boards, err := client.Boards()
	if err != nil {
		log.Fatal(err)
	}
	if len(boards) == 0 {
		fmt.Println("No boards found")
		return
	}
	for _, board := range boards {
		fmt.Println(board.Name)
	}
}
