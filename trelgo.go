package trelgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const API_PREFIX = "https://api.trello.com/1/"

type Client struct {
	Username string
	APIKey   string
	Token    string
}

type Board struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	client *Client
}

type List struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Closed  bool   `json:"closed"`
	IDBoard string `json:"idBoard"`
	Board   Board
	client  *Client
}

type Card struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Closed       bool     `json:"closed"`
	Description  string   `json:"desc"`
	IDBoard      string   `json:"idBoard"`
	IDChecklists []string `json:"idChecklists"`
	IDList       string   `json:"idList"`
	List         List
	Board        Board
	client       *Client
}

func New(username, apiKey, token string) *Client {
	return &Client{
		Username: username,
		APIKey:   apiKey,
		Token:    token,
	}
}

func (c *Client) Boards() ([]Board, error) {
	apiurl := API_PREFIX + fmt.Sprintf("members/%s/boards?key=%s&token=%s", c.Username, c.APIKey, c.Token)
	var out []Board
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].client = c
	}
	return out, nil
}

func (c *Client) Board(id string) (Board, error) {
	apiurl := API_PREFIX + fmt.Sprintf("boards/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Board
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return Board{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) Card(id string) (Card, error) {
	apiurl := API_PREFIX + fmt.Sprintf("cards/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Card
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return Card{}, err
	}
	out.client = c
	return out, nil
}

func (b Board) Lists() ([]List, error) {
	c := b.client
	apiurl := API_PREFIX + fmt.Sprintf("boards/%s/lists?key=%s&token=%s", b.ID, c.APIKey, c.Token)
	var out []List
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Board = b
		out[i].client = c
	}
	return out, nil
}

func (b Board) NewList(name, position string) (List, error) {
	c := b.client
	if position == "" {
		position = "bottom"
	}
	name, position = url.QueryEscape(name), url.QueryEscape(position)
	apiurl := API_PREFIX + fmt.Sprintf("boards/%s/lists?name=%s&pos=%s&key=%s&token=%s", b.ID, name, position, c.APIKey, c.Token)
	var out List
	if err := doMethodAndParseBody(http.MethodPost, apiurl, &out); err != nil {
		return List{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) List(id string) (List, error) {
	apiurl := API_PREFIX + fmt.Sprintf("lists/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out List
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return List{}, err
	}
	out.client = c
	return out, nil
}

func (l List) Cards() ([]Card, error) {
	c := l.client
	apiurl := API_PREFIX + fmt.Sprintf("lists/%s/cards?key=%s&token=%s", l.ID, c.APIKey, c.Token)
	var out []Card
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Board = l.Board
		out[i].List = l
		out[i].client = c
	}
	return out, nil
}

func (l List) NewCard(name, desc, position string) (Card, error) {
	c := l.client
	name, desc, position = url.QueryEscape(name), url.QueryEscape(desc), url.QueryEscape(position)
	query := fmt.Sprintf("idList=%s&name=%s&desc=%s&pos=%s&key=%s&token=%s", l.ID, name, desc, position, c.APIKey, c.Token)
	apiurl := API_PREFIX + "cards?" + query
	var out Card
	if err := doMethodAndParseBody(http.MethodPost, apiurl, &out); err != nil {
		return Card{}, err
	}
	out.Board = l.Board
	out.List = l
	out.client = c
	return out, nil
}

func (c *Card) Move(listID string) error {
	cl := c.client
	apiurl := API_PREFIX + fmt.Sprintf("cards/%s?idList=%s&key=%s&token=%s", c.ID, listID, cl.APIKey, cl.Token)
	if err := doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	// TODO: Eventually handle List and ListID mismatch in a better way.
	c.IDList = listID
	c.List.ID = listID
	return nil
}

func doMethod(method, apiurl string) error {
	req, err := http.NewRequest(method, apiurl, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close() // Not deferred because we ignore the body.
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP Request error, status: %d", resp.StatusCode)
	}
	return nil
}

// t must be a pointer
func doMethodAndParseBody(method, apiurl string, t interface{}) error {
	req, err := http.NewRequest(method, apiurl, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP Request error, status: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, t)
}
