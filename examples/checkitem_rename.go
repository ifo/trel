package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ifo/trel"
)

func main() {
	apiKey := flag.String("k", "apikey", "your trello api key")
	token := flag.String("t", "token", "your token - https://trello.com/app-key")
	checklistID := flag.String("ch", "checklist id", "your checklist")
	checkItemName := flag.String("ci", "check item name", "your check item's original name")
	checkItemNewName := flag.String("cinew", "new check item name", "your check item's new name")
	flag.Parse()

	client := trel.New(
		nil,     // Default http client
		*apiKey, // Your api key
		*token,  // Your token (https://trello.com/app-key)
	)

	checklist, err := client.Checklist(*checklistID)
	if err != nil {
		log.Fatalf("Error getting checklist: %s\n", err)
	}

	checkItem, err := checklist.CheckItems.Find(*checkItemName)
	if err != nil {
		log.Fatalf("Error getting checkItem: %s\n", err)
	}

	if err := checkItem.Rename(*checkItemNewName); err != nil {
		log.Fatalf("Error renaming checkItem: %s\n", err)
	}

	fmt.Printf("%v\n", checkItem)
}
