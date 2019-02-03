package trel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		t.Error(err)
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
			t.Error(err)
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
		t.Error(err)
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
		t.Error(err)
	}

	if checklist.Name != "Test" {
		t.Errorf("Expected %q, got %q\n", "Test", checklist.Name)
	}
	if checklist.ID != "1234" {
		t.Errorf("Expected %q, got %q\n", "1234", checklist.ID)
	}
	if checklist.IDBoard != "2345" {
		t.Errorf("Expected %q, got %q\n", "2345", checklist.IDBoard)
	}
	if checklist.IDCard != "3456" {
		t.Errorf("Expected %q, got %q\n", "3456", checklist.IDCard)
	}
	// Not checking checklist.Card or checklist.Board because
	// they are only set when a checklist is obtained from a board.
	checkItem := checklist.CheckItems[0]
	if checkItem.Name != "Test CheckItem" {
		t.Errorf("Expected %q, got %q\n", "Test CheckItem", checkItem.Name)
	}
	if checkItem.ID != "4567" {
		t.Errorf("Expected %q, got %q\n", "4567", checkItem.ID)
	}
	if checkItem.IDChecklist != "1234" {
		t.Errorf("Expected %q, got %q\n", "1234", checkItem.IDChecklist)
	}
	if checkItem.State != "incomplete" {
		t.Errorf("Expected %q, got %q\n", "incomplete", checkItem.State)
	}
}
