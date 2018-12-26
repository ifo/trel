package trelgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	url := API_PREFIX + fmt.Sprintf("members/%s/boards?key=%s&token=%s", c.Username, c.APIKey, c.Token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out []Board
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}
	for i := range out {
		out[i].client = c
	}
	return out, nil
}

func (c *Client) Board(id string) (Board, error) {
	url := API_PREFIX + fmt.Sprintf("boards/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	resp, err := http.Get(url)
	if err != nil {
		return Board{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Board{}, err
	}

	var out Board
	err = json.Unmarshal(body, &out)
	if err != nil {
		return Board{}, err
	}
	out.client = c
	return out, nil
}

func (b Board) Lists() ([]List, error) {
	c := b.client
	url := API_PREFIX + fmt.Sprintf("boards/%s/lists?key=%s&token=%s", b.ID, c.APIKey, c.Token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out []List
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Board = b
		out[i].client = c
	}
	return out, nil
}

func (c *Client) List(id string) (List, error) {
	url := API_PREFIX + fmt.Sprintf("lists/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	resp, err := http.Get(url)
	if err != nil {
		return List{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return List{}, err
	}

	var out List
	err = json.Unmarshal(body, &out)
	if err != nil {
		return List{}, err
	}
	out.client = c
	return out, nil
}

func (l List) Cards() ([]Card, error) {
	c := l.client
	url := API_PREFIX + fmt.Sprintf("lists/%s/cards?key=%s&token=%s", l.ID, c.APIKey, c.Token)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out []Card
	err = json.Unmarshal(body, &out)
	if err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Board = l.Board
		out[i].List = l
		out[i].client = c
	}
	return out, nil
}
