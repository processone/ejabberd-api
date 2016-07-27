package ejabberd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

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

//==============================================================================

type GetStats struct {
	Name string `json:"name"`
}

func (g *GetStats) knownStats() []string {
	return []string{"registeredusers", "onlineusers", "onlineusersnode",
		"uptimeseconds", "processes"}
}

func (g *GetStats) params() (HTTPParams, error) {
	var query url.Values
	if !stringInSlice(g.Name, g.knownStats()) {
		return HTTPParams{}, fmt.Errorf("unknow statistic: %s", g.Name)
	}

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

//==============================================================================

type RegisterUser struct {
	JID      string `json:"jid"`
	Password string `json:"password"`
}

func (r *RegisterUser) params() (HTTPParams, error) {
	var query url.Values

	jid, err := parseJID(r.JID)
	if err != nil {
		return HTTPParams{}, err
	}

	type register struct {
		User     string `json:"user"`
		Host     string `json:"host"`
		Password string `json:"password"`
	}

	data := register{
		User:     jid.username,
		Host:     jid.domain,
		Password: r.Password,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return HTTPParams{}, err
	}

	return HTTPParams{
		version: 1,
		admin:   true,
		method:  "POST",
		path:    "register/",
		query:   query,
		body:    body,
	}, nil
}
