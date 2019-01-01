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
	checklistID := flag.String("ch", "checklist id", "your checklist")
	flag.Parse()

	client := trel.New(
		*username, // Your username
		*apiKey,   // Your api key
		*token,    // Your token (https://trello.com/app-key)
	)

	checklist, err := client.Checklist(*checklistID)
	if err != nil {
		log.Fatal(err)
	}

	for i, checkitem := range checklist.CheckItems {
		err := checkitem.Complete()
		if err != nil {
			log.Fatal(err)
		}
		if i%2 == 0 {
			err = checkitem.Incomplete()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	for _, checkitem := range checklist.CheckItems {
		// Checkitems should alternate between "complete" and "incomplete"
		fmt.Println(checkitem.State)
	}
}
