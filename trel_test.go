package trel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setupClientMuxServer() (*Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := New(server.Client(), "", "")
	client.BaseURL, _ = url.Parse(server.URL)
	return client, mux, server
}

func TestClient_Boards(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"name": "Test", "id": "1234"}]`)
	})

	boards, err := client.Boards("Username")
	if err != nil {
		t.Fatal(err)
	}

	compare := Boards{{
		ID:     "1234",
		Name:   "Test",
		client: client,
	}}

	if !reflect.DeepEqual(compare, boards) {
		t.Errorf("Expected %#v, got %#v\n", compare, boards)
	}
}

func TestClient_Board(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		Board Board
		Body  string
	}{
		{Board: Board{Name: "Test", ID: "1234", client: client},
			Body: `{"name": "Test", "id": "1234"}`},
		{Board: Board{Name: "Board", ID: "5678", client: client},
			Body: `{"name": "Board", "id": "5678"}`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		board, err := client.Board(c.Board.ID)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.Board, board) {
			t.Errorf("Expected %#v, got %#v\n", c.Board, board)
		}
	}
}

func TestClient_List(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": "1234", "name": "Test", "idBoard": "5678"}`)
	})

	compareList := List{ID: "1234", Name: "Test", IDBoard: "5678", client: client}

	list, err := client.List("1234")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(compareList, list) {
		t.Errorf("Expected %#v, got %#v\n", compareList, list)
	}
}

func TestClient_Checklist(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": "1234", "name": "Test", "idBoard": "2345", "idCard": "3456", "checkItems": [
			{"idChecklist": "1234", "state": "incomplete", "id": "4567", "name": "Test CheckItem"}
		]}`)
	})

	checklist, err := client.Checklist("1234")
	if err != nil {
		t.Fatal(err)
	}

	compare := Checklist{
		ID:      "1234",
		Name:    "Test",
		IDBoard: "2345",
		IDCard:  "3456",
		Card:    Card{},
		Board:   Board{},
		CheckItems: CheckItems{{
			ID:          "4567",
			Name:        "Test CheckItem",
			IDChecklist: "1234",
			State:       "incomplete",
			client:      client,
			//Checklist: set this later,
		}},
		client: client,
	}
	// Properly set CheckItem's Checklist
	for i := range compare.CheckItems {
		compare.CheckItems[i].Checklist = compare
	}

	if !reflect.DeepEqual(compare, checklist) {
		t.Errorf("Expected %v, got %v\n", compare, checklist)
	}
}

func TestClient_NewWebhook(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		Webhook Webhook
		Body    string
	}{
		{Webhook: Webhook{ID: "1234", Description: "Card", IDModel: "5678",
			CallbackURL: "example.com/card/5678", Active: true, client: client},
			Body: `{"id": "1234", "description": "Card", "idModel": "5678", "callbackURL": "example.com/card/5678", "active": true}`},
		{Webhook: Webhook{ID: "2345", Description: "List", IDModel: "6789",
			CallbackURL: "example.com/list/6789", Active: true, client: client},
			Body: `{"id": "2345", "description": "List", "idModel": "6789", "callbackURL": "example.com/list/6789", "active": true}`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		webhook, err := client.NewWebhook(c.Webhook.Description, c.Webhook.CallbackURL, c.Webhook.IDModel)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.Webhook, webhook) {
			t.Errorf("Expected %#v, got %#v\n", c.Webhook, webhook)
		}
	}
}

func TestClient_Webhooks(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		Webhooks Webhooks
		Body     string
	}{
		{Webhooks: Webhooks{
			{ID: "1234", Description: "Card", IDModel: "5678",
				CallbackURL: "example.com/card/5678", Active: true, client: client},
			{ID: "2345", Description: "List", IDModel: "6789",
				CallbackURL: "example.com/list/6789", Active: true, client: client},
		},
			Body: `[{"id": "1234", "description": "Card", "idModel": "5678", "callbackURL": "example.com/card/5678", "active": true},
			{"id": "2345", "description": "List", "idModel": "6789", "callbackURL": "example.com/list/6789", "active": true}]`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		webhooks, err := client.Webhooks()
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.Webhooks, webhooks) {
			t.Errorf("Expected %#v, got %#v\n", c.Webhooks, webhooks)
		}
	}
}

func TestClient_Webhook(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		Webhook Webhook
		Body    string
	}{
		{Webhook: Webhook{ID: "1234", Description: "Card", IDModel: "5678",
			CallbackURL: "example.com/card/5678", Active: true, client: client},
			Body: `{"id": "1234", "description": "Card", "idModel": "5678", "callbackURL": "example.com/card/5678", "active": true}`},
		{Webhook: Webhook{ID: "2345", Description: "List", IDModel: "6789",
			CallbackURL: "example.com/list/6789", Active: true, client: client},
			Body: `{"id": "2345", "description": "List", "idModel": "6789", "callbackURL": "example.com/list/6789", "active": true}`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		webhook, err := client.Webhook(c.Webhook.ID)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.Webhook, webhook) {
			t.Errorf("Expected %#v, got %#v\n", c.Webhook, webhook)
		}
	}
}

func TestBoard_Lists(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	compareBoard := Board{ID: "1234", client: client}
	cases := []struct {
		Lists Lists
		Body  string
	}{
		{Lists: Lists{
			{ID: "2345", Name: "List 1", Closed: false, IDBoard: "1234", Board: compareBoard, client: client},
			{ID: "3456", Name: "List 2", Closed: false, IDBoard: "1234", Board: compareBoard, client: client},
		}, Body: `[{"id": "2345", "name": "List 1", "idBoard": "1234"}, {"id": "3456", "name": "List 2", "idBoard": "1234"}]`,
		},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		board := Board{ID: "1234", client: client}
		lists, err := board.Lists()
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(compareBoard, board) {
			t.Errorf("Expected %#v, got %#v\n", compareBoard, board)
		}

		if !reflect.DeepEqual(c.Lists, lists) {
			t.Errorf("Expected %#v, got %#v\n", c.Lists, lists)
		}
	}
}

func TestBoard_NewList(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	board := Board{ID: "1234", client: client}
	cases := []struct {
		List List
		Body string
	}{
		{List: List{ID: "2345", Name: "List 1", Closed: false, IDBoard: "1234", Board: board, client: client},
			Body: `{"id": "2345", "name": "List 1", "idBoard": "1234"}`},
		{List: List{ID: "3456", Name: "List 2", Closed: false, IDBoard: "1234", Board: board, client: client},
			Body: `{"id": "3456", "name": "List 2", "idBoard": "1234"}`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		list, err := board.NewList(c.List.Name, "")
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(c.List, list) {
			t.Errorf("Expected %#v, got %#v\n", c.List, list)
		}
	}
}

func TestBoard_FindList(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	board := Board{ID: "1234", client: client}
	cases := []struct {
		ListName string
		List     List
		Body     string
		Err      error
	}{
		{ListName: "List 1",
			List: List{ID: "2345", Name: "List 1", Closed: false, IDBoard: "1234", Board: board, client: client},
			Body: `[{"id": "2345", "name": "List 1", "idBoard": "1234"}, {"id": "3456", "name": "List 2", "idBoard": "1234"}]`,
			Err:  nil},
		{ListName: "List 2",
			List: List{ID: "3456", Name: "List 2", Closed: false, IDBoard: "1234", Board: board, client: client},
			Body: `[{"id": "2345", "name": "List 1", "idBoard": "1234"}, {"id": "3456", "name": "List 2", "idBoard": "1234"}]`,
			Err:  nil},
		{ListName: "List 1",
			List: List{},
			Body: `[{"id": "3456", "name": "List 2", "idBoard": "1234"}]`,
			Err:  NotFoundError{Type: "List", Identifier: "List 1"}},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		list, err := board.FindList(c.ListName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.List, list) {
			t.Errorf("Expected %#v, got %#v\n", c.List, list)
		}
	}
}

func TestList_Cards(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	list := List{ID: "1234", client: client}
	cases := []struct {
		Cards Cards
		Body  string
	}{
		{Cards: Cards{
			{ID: "2345", Name: "Card 1", IDList: "1234", List: list, client: client},
			{ID: "3456", Name: "Card 2", IDList: "1234", List: list, client: client}},
			Body: `[{"id": "2345", "name": "Card 1", "idList": "1234"}, {"id": "3456", "name": "Card 2", "idList": "1234"}]`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		cards, err := list.Cards()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.Cards, cards) {
			t.Errorf("Expected %#v, got %#v\n", c.Cards, cards)
		}
	}
}

func TestList_FindCard(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	list := List{ID: "1234", client: client}
	cases := []struct {
		CardName string
		Card     Card
		Body     string
		Err      error
	}{
		{CardName: "Card 1",
			Card: Card{ID: "2345", Name: "Card 1", IDList: "1234", List: list, client: client},
			Body: `[{"id": "2345", "name": "Card 1", "idList": "1234"}, {"id": "3456", "name": "Card 2", "idList": "1234"}]`,
			Err:  nil},
		{CardName: "Card 2",
			Card: Card{ID: "3456", Name: "Card 2", IDList: "1234", List: list, client: client},
			Body: `[{"id": "2345", "name": "Card 1", "idList": "1234"}, {"id": "3456", "name": "Card 2", "idList": "1234"}]`,
			Err:  nil},
		{CardName: "Card 1",
			Card: Card{},
			Body: `[{"id": "3456", "name": "Card 2", "idList": "1234"}]`,
			Err:  NotFoundError{Type: "Card", Identifier: "Card 1"}},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		card, err := list.FindCard(c.CardName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.Card, card) {
			t.Errorf("Expected %#v, got %#v\n", c.Card, card)
		}
	}
}

func TestList_NewCard(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	board := Board{ID: "4321"}
	list := List{ID: "1234", Board: board, client: client}
	cases := []struct {
		Card Card
		Body string
	}{
		{Card: Card{ID: "2345", Name: "Card 1", Description: "first card",
			IDList: list.ID, List: list, IDBoard: board.ID, Board: board, client: client},
			Body: `{"id": "2345", "name": "Card 1", "desc": "first card", "idList": "1234", "idBoard": "4321"}`},
		{Card: Card{ID: "3456", Name: "Card 2", Description: "second card",
			IDList: list.ID, List: list, IDBoard: board.ID, Board: board, client: client},
			Body: `{"id": "3456", "name": "Card 2", "desc": "second card", "idList": "1234", "idBoard": "4321"}`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		card, err := list.NewCard(c.Card.Name, c.Card.Description, "")
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.Card, card) {
			t.Errorf("Expected %#v, got %#v\n", c.Card, card)
		}
	}
}

func TestLists_Find(t *testing.T) {
	list1 := List{ID: "2345", Name: "List 1"}
	list2 := List{ID: "3456", Name: "List 2"}
	list3 := List{ID: "4567", Name: "List 3"}
	cases := []struct {
		ListName  string
		FoundList *List
		Lists     Lists
		Err       error
	}{
		{ListName: list1.Name,
			FoundList: &list1,
			Lists:     Lists{list1, list2, list3},
			Err:       nil},
		{ListName: list2.Name,
			FoundList: &list2,
			Lists:     Lists{list1, list2, list3},
			Err:       nil},
		{ListName: list3.Name,
			FoundList: &list3,
			Lists:     Lists{list1, list2, list3},
			Err:       nil},
		{ListName: list1.Name,
			FoundList: &List{},
			Lists:     Lists{},
			Err:       NotFoundError{Type: "List", Identifier: list1.Name}},
		{ListName: list2.Name,
			FoundList: &List{},
			Lists:     Lists{list1, list3},
			Err:       NotFoundError{Type: "List", Identifier: list2.Name}},
	}

	for _, c := range cases {
		list, err := c.Lists.Find(c.ListName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.FoundList, list) {
			t.Errorf("Expected %#v, got %#v\n", c.FoundList, list)
		}
	}
}

func TestCard_Move(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	list1 := List{ID: "2345"}
	list2 := List{ID: "3456"}
	cases := []struct {
		ListID  string
		Card    Card
		EndCard Card
		Err     error
	}{
		{ListID: list2.ID,
			Card:    Card{IDList: list1.ID, List: list1, client: client},
			EndCard: Card{IDList: list2.ID, List: list2, client: client},
			Err:     nil},
		{ListID: list1.ID,
			Card:    Card{IDList: list2.ID, List: list2, client: client},
			EndCard: Card{IDList: list1.ID, List: list1, client: client},
			Err:     nil},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	for _, c := range cases {
		err := c.Card.Move(c.ListID)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.Card, c.EndCard) {
			t.Errorf("Expected %#v, got %#v\n", c.Card, c.EndCard)
		}
	}
}

func TestCard_Rename(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	name1 := "card 1"
	name2 := "card 2"
	cases := []struct {
		NewName string
		Card    Card
		EndCard Card
		Err     error
	}{
		{NewName: name2,
			Card:    Card{Name: name1, client: client},
			EndCard: Card{Name: name2, client: client},
			Err:     nil},
		{NewName: name1,
			Card:    Card{Name: name2, client: client},
			EndCard: Card{Name: name1, client: client},
			Err:     nil},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	for _, c := range cases {
		err := c.Card.Rename(c.NewName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.Card, c.EndCard) {
			t.Errorf("Expected %#v, got %#v\n", c.Card, c.EndCard)
		}
	}
}

func TestCard_Checklists(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	card := Card{client: client}
	cases := []struct {
		Checklists Checklists
		Body       string
	}{
		{Checklists: Checklists{{ID: "1234", Name: "Checklist 1", Card: card, client: client, CheckItems: nil}},
			Body: `[{"id": "1234", "name": "Checklist 1"}]`},
		{Checklists: Checklists{{ID: "1234", Name: "Checklist 1", Card: card, client: client, CheckItems: CheckItems{
			{ID: "2345", Name: "CheckItem 1", State: "incomplete", IDChecklist: "1234", client: client}}}},
			Body: `[{"id": "1234", "name": "Checklist 1", "checkItems": [
				{"idChecklist": "1234", "state": "incomplete", "id": "2345", "name": "CheckItem 1"}]}]`},
	}
	// Properly set the Checklist on each CheckItem.
	for _, c := range cases {
		for _, checklist := range c.Checklists {
			for i := range checklist.CheckItems {
				checklist.CheckItems[i].Checklist = checklist
			}
		}
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		checklists, err := card.Checklists()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.Checklists, checklists) {
			t.Errorf("Expected %v, got %v\n", c.Checklists, checklists)
		}
	}
}

func TestCards_Find(t *testing.T) {
	card1 := Card{ID: "2345", Name: "Card 1"}
	card2 := Card{ID: "3456", Name: "Card 2"}
	card3 := Card{ID: "4567", Name: "Card 3"}
	cases := []struct {
		CardName  string
		FoundCard *Card
		Cards     Cards
		Err       error
	}{
		{CardName: card1.Name,
			FoundCard: &card1,
			Cards:     Cards{card1, card2, card3},
			Err:       nil},
		{CardName: card2.Name,
			FoundCard: &card2,
			Cards:     Cards{card1, card2, card3},
			Err:       nil},
		{CardName: card3.Name,
			FoundCard: &card3,
			Cards:     Cards{card1, card2, card3},
			Err:       nil},
		{CardName: card1.Name,
			FoundCard: &Card{},
			Cards:     Cards{},
			Err:       NotFoundError{Type: "Card", Identifier: card1.Name}},
		{CardName: card2.Name,
			FoundCard: &Card{},
			Cards:     Cards{card1, card3},
			Err:       NotFoundError{Type: "Card", Identifier: card2.Name}},
	}

	for _, c := range cases {
		card, err := c.Cards.Find(c.CardName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.FoundCard, card) {
			t.Errorf("Expected %#v, got %#v\n", c.FoundCard, card)
		}
	}
}

func TestCheckItem_Complete(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		CheckItem    CheckItem
		EndCheckItem CheckItem
	}{
		{CheckItem: CheckItem{State: "incomplete", client: client},
			EndCheckItem: CheckItem{State: "complete", client: client}},
		{CheckItem: CheckItem{State: "complete", client: client},
			EndCheckItem: CheckItem{State: "complete", client: client}},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	for _, c := range cases {
		err := c.CheckItem.Complete()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.CheckItem, c.EndCheckItem) {
			t.Errorf("Expected %#v, got %#v\n", c.CheckItem, c.EndCheckItem)
		}
	}
}

func TestCheckItem_Incomplete(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		CheckItem    CheckItem
		EndCheckItem CheckItem
	}{
		{CheckItem: CheckItem{State: "complete", client: client},
			EndCheckItem: CheckItem{State: "incomplete", client: client}},
		{CheckItem: CheckItem{State: "incomplete", client: client},
			EndCheckItem: CheckItem{State: "incomplete", client: client}},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	for _, c := range cases {
		err := c.CheckItem.Incomplete()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(c.CheckItem, c.EndCheckItem) {
			t.Errorf("Expected %#v, got %#v\n", c.CheckItem, c.EndCheckItem)
		}
	}
}

func TestCheckItems_Find(t *testing.T) {
	checkitem1 := CheckItem{ID: "2345", Name: "CheckItem 1"}
	checkitem2 := CheckItem{ID: "3456", Name: "CheckItem 2"}
	checkitem3 := CheckItem{ID: "4567", Name: "CheckItem 3"}
	cases := []struct {
		CheckItemName  string
		FoundCheckItem *CheckItem
		CheckItems     CheckItems
		Err            error
	}{
		{CheckItemName: checkitem1.Name,
			FoundCheckItem: &checkitem1,
			CheckItems:     CheckItems{checkitem1, checkitem2, checkitem3},
			Err:            nil},
		{CheckItemName: checkitem2.Name,
			FoundCheckItem: &checkitem2,
			CheckItems:     CheckItems{checkitem1, checkitem2, checkitem3},
			Err:            nil},
		{CheckItemName: checkitem3.Name,
			FoundCheckItem: &checkitem3,
			CheckItems:     CheckItems{checkitem1, checkitem2, checkitem3},
			Err:            nil},
		{CheckItemName: checkitem1.Name,
			FoundCheckItem: &CheckItem{},
			CheckItems:     CheckItems{},
			Err:            NotFoundError{Type: "CheckItem", Identifier: checkitem1.Name}},
		{CheckItemName: checkitem2.Name,
			FoundCheckItem: &CheckItem{},
			CheckItems:     CheckItems{checkitem1, checkitem3},
			Err:            NotFoundError{Type: "CheckItem", Identifier: checkitem2.Name}},
	}

	for _, c := range cases {
		checkitem, err := c.CheckItems.Find(c.CheckItemName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.FoundCheckItem, checkitem) {
			t.Errorf("Expected %#v, got %#v\n", c.FoundCheckItem, checkitem)
		}
	}
}

func TestCheckItem_Rename(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	name1 := "checkitem 1"
	name2 := "checkitem 2"
	cases := []struct {
		NewName      string
		CheckItem    CheckItem
		EndCheckItem CheckItem
		Err          error
	}{
		{NewName: name2,
			CheckItem:    CheckItem{Name: name1, client: client},
			EndCheckItem: CheckItem{Name: name2, client: client},
			Err:          nil},
		{NewName: name1,
			CheckItem:    CheckItem{Name: name2, client: client},
			EndCheckItem: CheckItem{Name: name1, client: client},
			Err:          nil},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	})

	for _, c := range cases {
		err := c.CheckItem.Rename(c.NewName)
		if c.Err != err {
			t.Errorf("Expected %q, got %q\n", c.Err, err)
		}

		if !reflect.DeepEqual(c.CheckItem, c.EndCheckItem) {
			t.Errorf("Expected %#v, got %#v\n", c.CheckItem, c.EndCheckItem)
		}
	}
}
