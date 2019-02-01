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
		fmt.Fprint(w, `[{"name": "Test", "id": "111122223333444455556666"}]`)
	})

	boards, err := client.Boards("Username")
	if err != nil {
		t.Error(err)
	}

	if len(boards) != 1 {
		t.Errorf("Expected 1 Board, got %d boards\n", len(boards))
	}

	if boards[0].Name != "Test" {
		t.Errorf("Expected \"Test\", got %q\n", boards[0].Name)
	}
	if boards[0].ID != "111122223333444455556666" {
		t.Errorf("Expected %q, got %q\n", "111122223333444455556666", boards[0].ID)
	}
}

func TestClient_Board(t *testing.T) {
	client, mux, server := setupClientMuxServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"name": "Test", "id": "111122223333444455556666"}`)
	})

	board, err := client.Board("111122223333444455556666")
	if err != nil {
		t.Error(err)
	}

	if board.Name != "Test" {
		t.Errorf("Expected \"Test\", got %q\n", board.Name)
	}
	if board.ID != "111122223333444455556666" {
		t.Errorf("Expected %q, got %q\n", "111122223333444455556666", board.ID)
	}
}
