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

type Checklist struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	IDBoard    string      `json:"idBoard"`
	IDCard     string      `json:"idCard"`
	CheckItems []CheckItem `json:"checkItems"`
	Card       Card
	Board      Board
	client     *Client
}

type CheckItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"` // TODO: Turn this into a boolean type and add custom json parsing.
	IDChecklist string `json:"idChecklist"`
	Checklist   Checklist
	client      *Client
}

type Webhook struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	IDModel     string `json:"idModel"`
	CallbackURL string `json:"callbackURL"` // TODO: Make this a url.URL instead of a string; add custom json parsing.
	Active      bool   `json:"active"`
	client      *Client
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

func (c *Client) Checklist(id string) (Checklist, error) {
	apiurl := API_PREFIX + fmt.Sprintf("checklists/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Checklist
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return Checklist{}, err
	}
	out.client = c
	for i := range out.CheckItems {
		out.CheckItems[i].Checklist = out
		out.CheckItems[i].client = c
	}
	return out, nil
}

func (c *Client) NewWebhook(description, callbackURL, idModel string) (Webhook, error) {
	description, callbackURL = url.QueryEscape(description), url.QueryEscape(callbackURL)
	apiurl := API_PREFIX + fmt.Sprintf("webhooks/?description=%s&callbackURL=%s&idModel=%s&key=%s&token=%s", description, callbackURL, idModel, c.APIKey, c.Token)
	var out Webhook
	if err := doMethodAndParseBody(http.MethodPost, apiurl, &out); err != nil {
		return Webhook{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) Webhooks() ([]Webhook, error) {
	apiurl := API_PREFIX + fmt.Sprintf("tokens/%s/webhooks?key=%s", c.Token, c.APIKey)
	var out []Webhook
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].client = c
	}
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
	// TODO: Handle query arguments and escaping better, probably use url.URL.
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

func (l List) FindCard(name string) (Card, error) {
	cards, err := l.Cards()
	if err != nil {
		return Card{}, err
	}
	for _, card := range cards {
		if card.Name == name {
			return card, nil
		}
	}
	return Card{}, fmt.Errorf("Card not found with name: %s", name)
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

func (ca *Card) Move(listID string) error {
	c := ca.client
	apiurl := API_PREFIX + fmt.Sprintf("cards/%s?idList=%s&key=%s&token=%s", ca.ID, listID, c.APIKey, c.Token)
	if err := doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	// TODO: Eventually handle List and ListID mismatch in a better way.
	ca.IDList = listID
	ca.List.ID = listID
	return nil
}

func (ca *Card) Checklists() ([]Checklist, error) {
	c := ca.client
	apiurl := API_PREFIX + fmt.Sprintf("cards/%s/checklists?key=%s&token=%s", ca.ID, c.APIKey, c.Token)
	var out []Checklist
	if err := doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].Card = *ca
		out[i].Board = ca.Board
		out[i].client = c
		for j := range out[i].CheckItems {
			// Properly set the Checklist for every CheckItem.
			out[i].CheckItems[j].Checklist = out[i]
			out[i].CheckItems[j].client = c
		}
	}
	return out, nil
}

func (ci *CheckItem) Complete() error {
	c := ci.client
	apiurl := API_PREFIX + fmt.Sprintf("cards/%s/checkItem/%s?state=complete&key=%s&token=%s", ci.Checklist.IDCard, ci.ID, c.APIKey, c.Token)
	if err := doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	ci.State = "complete"
	return nil
}

func (ci *CheckItem) Incomplete() error {
	c := ci.client
	apiurl := API_PREFIX + fmt.Sprintf("cards/%s/checkItem/%s?state=incomplete&key=%s&token=%s", ci.Checklist.IDCard, ci.ID, c.APIKey, c.Token)
	if err := doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	ci.State = "incomplete"
	return nil
}

func (w *Webhook) Activate() error {
	c := w.client
	apiurl := API_PREFIX + fmt.Sprintf("webhooks/%s?active=true&key=%s&token=%s", w.ID, c.APIKey, c.Token)
	if err := doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	w.Active = true
	return nil
}

func (w *Webhook) Deactivate() error {
	c := w.client
	apiurl := API_PREFIX + fmt.Sprintf("webhooks/%s?active=false&key=%s&token=%s", w.ID, c.APIKey, c.Token)
	if err := doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	w.Active = false
	return nil
}

func (w *Webhook) Delete() error {
	c := w.client
	apiurl := API_PREFIX + fmt.Sprintf("webhooks/%s?key=%s&token=%s", w.ID, c.APIKey, c.Token)
	if err := doMethod(http.MethodDelete, apiurl); err != nil {
		return err
	}
	w = &Webhook{}
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

// t must be a pointer.
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
