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

	if len(boards) != 1 {
		t.Errorf("Expected 1 Board, got %d boards\n", len(boards))
	}

	if boards[0].Name != "Test" {
		t.Errorf("Expected %q, got %q\n", "Test", boards[0].Name)
	}
	if boards[0].ID != "1234" {
		t.Errorf("Expected %q, got %q\n", "1234", boards[0].ID)
	}
}

func TestClient_Board(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	cases := []struct {
		Name string
		ID   string
		Body string
	}{
		{Name: "Test", ID: "1234", Body: `{"name": "Test", "id": "1234"}`},
		{Name: "Board", ID: "1234", Body: `{"name": "Board", "id": "1234"}`},
	}

	body := ""
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	})

	for _, c := range cases {
		body = c.Body

		board, err := client.Board(c.ID)
		if err != nil {
			t.Fatal(err)
		}

		if board.Name != c.Name {
			t.Errorf("Expected %q, got %q\n", c.Name, board.Name)
		}
		if board.ID != "1234" {
			t.Errorf("Expected %q, got %q\n", c.ID, board.ID)
		}
	}
}

func TestClient_List(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"id": "1234", "name": "Test", "closed": false, "idBoard": "5678"}`)
	})

	list, err := client.List("1234")
	if err != nil {
		t.Fatal(err)
	}

	if list.Name != "Test" {
		t.Errorf("Expected %q, got %q\n", "Test", list.Name)
	}
	if list.ID != "1234" {
		t.Errorf("Expected %q, got %q\n", "1234", list.ID)
	}
	if list.IDBoard != "5678" {
		t.Errorf("Expected %q, got %q\n", "5678", list.IDBoard)
	}
	if list.Closed != false {
		t.Errorf("Expected %t, got %t\n", false, list.Closed)
	}
	// Not checking list.Board because it is only set when a list is obtained from a board.
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
	// Properly set CheckItem Checklist
	for i := range compare.CheckItems {
		compare.CheckItems[i].Checklist = compare
	}

	if !reflect.DeepEqual(compare, checklist) {
		t.Errorf("Expected %v\n\nGot %v\n", compare, checklist)
	}
}
