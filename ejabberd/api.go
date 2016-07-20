package ejabberd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client is an ejabberd client API wrapper.
type Client struct {
	URL   string
	Token string
}

// Call performs the HTTP call to the API given client parameters.
func (c Client) Call(comm command) ([]byte, error) {
	p, err := comm.params()
	if err != nil {
		return []byte{}, err
	}

	url := c.URL + p.path
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(p.body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	if p.admin {
		req.Header.Set("X-Admin", "true")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// HTTPParams gathers all values needed by the client to encode actual
// ejabberd API call.
type HTTPParams struct {
	version int
	admin   bool
	method  string
	path    string
	query   url.Values
	body    []byte
}

type command interface {
	params() (HTTPParams, error)
}

type GetStats struct {
	Name string `json:"name"`
}

func (g *GetStats) params() (HTTPParams, error) {
	var query url.Values
	body, err := json.Marshal(g)
	if err != nil {
		return HTTPParams{}, err
	}

	return HTTPParams{
		version: 1,
		admin:   true,
		method:  "POST",
		path:    "stats/",
		query:   query,
		body:    body,
	}, nil
}
