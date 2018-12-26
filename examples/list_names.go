package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ifo/trelgo"
)

func main() {
	username := flag.String("u", "username", "your trello username or id")
	apiKey := flag.String("k", "apikey", "your trello api key")
	token := flag.String("t", "token", "your token - https://trello.com/app-key")
	boardID := flag.String("b", "board id", "your board's id")
	flag.Parse()

	client := trelgo.New(
		*username, // Your username
		*apiKey,   // Your api key
		*token,    // Your token (https://trello.com/app-key)
	)

	board, err := client.Board(*boardID)
	if err != nil {
		log.Fatal(err)
	}

	lists, err := board.Lists()
	if err != nil {
		log.Fatal(err)
	}
	for _, list := range lists {
		fmt.Println(list.Name)
	}
}
