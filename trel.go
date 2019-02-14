package trel

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const defaultAPIPrefix = "https://api.trello.com/1/"

type Client struct {
	client *http.Client

	BaseURL *url.URL

	APIKey string
	Token  string
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
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	IDBoard    string     `json:"idBoard"`
	IDCard     string     `json:"idCard"`
	CheckItems CheckItems `json:"checkItems"`
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

type Boards []Board
type Lists []List
type Cards []Card
type Checklists []Checklist
type CheckItems []CheckItem
type Webhooks []Webhook

func New(client *http.Client, apiKey, token string) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultAPIPrefix)

	return &Client{
		client:  client,
		BaseURL: baseURL,
		APIKey:  apiKey,
		Token:   token,
	}
}

func (c *Client) Boards(username string) (Boards, error) {
	apiurl := fmt.Sprintf("members/%s/boards?key=%s&token=%s", username, c.APIKey, c.Token)
	var out Boards
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].client = c
	}
	return out, nil
}

func (c *Client) Board(id string) (Board, error) {
	apiurl := fmt.Sprintf("boards/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Board
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return Board{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) List(id string) (List, error) {
	apiurl := fmt.Sprintf("lists/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out List
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return List{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) Card(id string) (Card, error) {
	apiurl := fmt.Sprintf("cards/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Card
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return Card{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) Checklist(id string) (Checklist, error) {
	apiurl := fmt.Sprintf("checklists/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Checklist
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
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
	apiurl := fmt.Sprintf("webhooks/?description=%s&callbackURL=%s&idModel=%s&key=%s&token=%s", description, callbackURL, idModel, c.APIKey, c.Token)
	var out Webhook
	if err := c.doMethodAndParseBody(http.MethodPost, apiurl, &out); err != nil {
		return Webhook{}, err
	}
	out.client = c
	return out, nil
}

func (c *Client) Webhooks() (Webhooks, error) {
	apiurl := fmt.Sprintf("tokens/%s/webhooks?key=%s", c.Token, c.APIKey)
	var out Webhooks
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return nil, err
	}
	for i := range out {
		out[i].client = c
	}
	return out, nil
}

func (c *Client) Webhook(id string) (Webhook, error) {
	apiurl := fmt.Sprintf("webhooks/%s?key=%s&token=%s", id, c.APIKey, c.Token)
	var out Webhook
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
		return Webhook{}, err
	}
	out.client = c
	return out, nil
}

func (b Board) Lists() (Lists, error) {
	c := b.client
	apiurl := fmt.Sprintf("boards/%s/lists?key=%s&token=%s", b.ID, c.APIKey, c.Token)
	var out Lists
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
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
	apiurl := fmt.Sprintf("boards/%s/lists?name=%s&pos=%s&key=%s&token=%s", b.ID, name, position, c.APIKey, c.Token)
	var out List
	if err := c.doMethodAndParseBody(http.MethodPost, apiurl, &out); err != nil {
		return List{}, err
	}
	out.client = c
	out.Board = b
	return out, nil
}

func (b Board) FindList(name string) (List, error) {
	lists, err := b.Lists()
	if err != nil {
		return List{}, err
	}
	l, err := lists.Find(name)
	return *l, err
}

func (l List) Cards() (Cards, error) {
	c := l.client
	apiurl := fmt.Sprintf("lists/%s/cards?key=%s&token=%s", l.ID, c.APIKey, c.Token)
	var out Cards
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
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
	c, err := cards.Find(name)
	return *c, err
}

func (l List) NewCard(name, desc, position string) (Card, error) {
	c := l.client
	name, desc, position = url.QueryEscape(name), url.QueryEscape(desc), url.QueryEscape(position)
	query := fmt.Sprintf("idList=%s&name=%s&desc=%s&pos=%s&key=%s&token=%s", l.ID, name, desc, position, c.APIKey, c.Token)
	apiurl := "cards?" + query
	var out Card
	if err := c.doMethodAndParseBody(http.MethodPost, apiurl, &out); err != nil {
		return Card{}, err
	}
	out.Board = l.Board
	out.List = l
	out.client = c
	return out, nil
}

func (ls Lists) Find(name string) (*List, error) {
	for i := range ls {
		if ls[i].Name == name {
			return &ls[i], nil
		}
	}
	return &List{}, NotFoundError{Type: "List", Identifier: name}
}

func (ca *Card) Move(listID string) error {
	// Don't do anything if the card is already on that list.
	if ca.IDList == listID {
		return nil
	}

	c := ca.client
	apiurl := fmt.Sprintf("cards/%s?idList=%s&key=%s&token=%s", ca.ID, listID, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	// TODO: Eventually handle List and ListID mismatch in a better way.
	ca.IDList = listID
	ca.List.ID = listID
	return nil
}

func (ca *Card) Rename(name string) error {
	if ca.Name == name {
		return nil
	}

	c := ca.client
	escapedName := url.QueryEscape(name)
	apiurl := fmt.Sprintf("cards/%s?name=%s&key=%s&token=%s", ca.ID, escapedName, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	ca.Name = name
	return nil
}

func (ca *Card) Checklists() (Checklists, error) {
	c := ca.client
	apiurl := fmt.Sprintf("cards/%s/checklists?key=%s&token=%s", ca.ID, c.APIKey, c.Token)
	var out Checklists
	if err := c.doMethodAndParseBody(http.MethodGet, apiurl, &out); err != nil {
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

func (cs Cards) Find(name string) (*Card, error) {
	for i := range cs {
		if cs[i].Name == name {
			return &cs[i], nil
		}
	}
	return &Card{}, NotFoundError{Type: "Card", Identifier: name}
}

func (ci *CheckItem) Complete() error {
	c := ci.client
	apiurl := fmt.Sprintf("cards/%s/checkItem/%s?state=complete&key=%s&token=%s", ci.Checklist.IDCard, ci.ID, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	ci.State = "complete"
	return nil
}

func (ci *CheckItem) Incomplete() error {
	c := ci.client
	apiurl := fmt.Sprintf("cards/%s/checkItem/%s?state=incomplete&key=%s&token=%s", ci.Checklist.IDCard, ci.ID, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	ci.State = "incomplete"
	return nil
}

func (cis CheckItems) Find(name string) (*CheckItem, error) {
	for i := range cis {
		if cis[i].Name == name {
			return &cis[i], nil
		}
	}
	return &CheckItem{}, NotFoundError{Type: "CheckItem", Identifier: name}
}

func (w *Webhook) Activate() error {
	// Don't activate active webhooks.
	if w.Active {
		return nil
	}

	c := w.client
	apiurl := fmt.Sprintf("webhooks/%s?active=true&key=%s&token=%s", w.ID, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	w.Active = true
	return nil
}

func (w *Webhook) Deactivate() error {
	// Don't deactivate inactive webhooks.
	if !w.Active {
		return nil
	}

	c := w.client
	apiurl := fmt.Sprintf("webhooks/%s?active=false&key=%s&token=%s", w.ID, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodPut, apiurl); err != nil {
		return err
	}
	w.Active = false
	return nil
}

func (w *Webhook) Delete() error {
	c := w.client
	apiurl := fmt.Sprintf("webhooks/%s?key=%s&token=%s", w.ID, c.APIKey, c.Token)
	if err := c.doMethod(http.MethodDelete, apiurl); err != nil {
		return err
	}
	w = &Webhook{}
	return nil
}

func (ws Webhooks) Find(modelID string) (*Webhook, error) {
	for i := range ws {
		if ws[i].IDModel == modelID {
			return &ws[i], nil
		}
	}
	return &Webhook{}, NotFoundError{Type: "Webhook", Identifier: modelID}
}

func (c *Client) doMethod(method, apiurl string) error {
	reqURL := joinPath(c.BaseURL.String(), apiurl)
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close() // Not deferred because we ignore the body.
	if resp.StatusCode >= 400 {
		return HTTPRequestError{StatusCode: resp.StatusCode}
	}
	return nil
}

// t must be a pointer.
func (c *Client) doMethodAndParseBody(method, apiurl string, t interface{}) error {
	reqURL := joinPath(c.BaseURL.String(), apiurl)
	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return HTTPRequestError{StatusCode: resp.StatusCode}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, t)
}

type NotFoundError struct {
	Type       string
	Identifier string
}

func (n NotFoundError) Error() string {
	return fmt.Sprintf("%s with identifier %q was not found", n.Type, n.Identifier)
}

type HTTPRequestError struct {
	StatusCode int
}

func (h HTTPRequestError) Error() string {
	return fmt.Sprintf("HTTP Request error with status: %d", h.StatusCode)
}

func joinPath(host, path string) string {
	runeHost, runePath := []rune(host), []rune(path)
	if runeHost[len(runeHost)-1] != '/' && runePath[0] != '/' {
		return host + "/" + path
	}
	return host + path
}

func (cl Checklist) String() string {
	return fmt.Sprintf("Checklist -\nID: %q,\nName: %q,\nIDBoard: %q,\nIDCard: %q,\nBoard: %v,\nCard: %v,\nclient: %v,\nCheckItems:\n%s",
		cl.ID,
		cl.Name,
		cl.IDBoard,
		cl.IDCard,
		cl.Board,
		cl.Card,
		cl.client,
		cl.CheckItems.String(),
	)
}

func (ci CheckItem) String() string {
	return fmt.Sprintf("CheckItem - ID: %q, Name: %q, State: %q, IDChecklist: %q, Checklist: %q, client: %v",
		ci.ID,
		ci.Name,
		ci.State,
		ci.IDChecklist,
		ci.Checklist.Name,
		ci.client,
	)
}

func (cis CheckItems) String() string {
	out := ""
	for _, ci := range cis {
		out += ci.String() + "\n"
	}
	return out
}
