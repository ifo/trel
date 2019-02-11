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
