package ejabberd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPParams gathers all values needed by the client to encode actual
// ejabberd API call.
type APIParams struct {
	name    string
	version int
	admin   bool // Does API require admin header ?

	method string
	query  url.Values
	body   []byte
}

// Request is the common interface to all ejabberd requests.
type Request interface {
	params() (APIParams, error)
	parseResponse([]byte) (Response, error)
}

// Response is the command interface for all ejabberd API call
// results.
type Response interface {
	errorCode() int
}

// Call performs the HTTP call to ejabberd API given client
// parameters. It returns a struct complying with Response interface.
func (c Client) Call(req Request) (Response, error) {
	resp, err := c.CallRaw(req)
	if err != nil {
		return nil, err
	}

	return req.parseResponse(resp)
}

// CallRaw performs HTTP call to ejabberd API and returns Raw Body
// reponse from the server as slice of bytes.
func (c Client) CallRaw(req Request) ([]byte, error) {
	p, err := req.params()
	if err != nil {
		return []byte{}, err
	}

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}

	url := c.BaseURL + p.name + "/"
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(p.body))
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	if p.admin {
		r.Header.Set("X-Admin", "true")
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

//==============================================================================

type StatsRequest struct {
	Name string `json:"name"`
}

func (g StatsRequest) knownStats() []string {
	return []string{"registeredusers", "onlineusers", "onlineusersnode",
		"uptimeseconds", "processes"}
}

func (g *StatsRequest) params() (APIParams, error) {
	var query url.Values
	if !stringInSlice(g.Name, g.knownStats()) {
		return APIParams{}, fmt.Errorf("unknow statistic: %s", g.Name)
	}

	body, err := json.Marshal(g)
	if err != nil {
		return APIParams{}, err
	}

	return APIParams{
		name:    "stats",
		version: 1,

		admin:  true,
		method: "POST",

		query: query,
		body:  body,
	}, nil
}

//==============================================================================

type RegisterRequest struct {
	JID      string `json:"jid"`
	Password string `json:"password"`
}

func (r *RegisterRequest) params() (APIParams, error) {
	var query url.Values

	jid, err := parseJID(r.JID)
	if err != nil {
		return APIParams{}, err
	}

	// Actual parameter for ejabberd. We expose JID string as it is
	// easier to manipulate from a client.
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
		return APIParams{}, err
	}

	return APIParams{
		name:    "register",
		version: 1,
		admin:   true,

		method: "POST",
		query:  query,
		body:   body,
	}, nil
}
