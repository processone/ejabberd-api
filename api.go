package ejabberd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Call performs HTTP call to ejabberd API given client parameters. It
// returns a struct complying with Response interface.
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

	// TODO: We should limit the amount of data the client reads from ejabberd as response
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return body, err
}

// Request is the common interface to all ejabberd requests.
type Request interface {
	params() (apiParams, error)
	parseResponse([]byte) (Response, error)
}

// Response is the command interface for all ejabberd API call
// results.
type Response interface {
	JSON() string
}

// apiParams gathers all values needed by the client to encode actual
// ejabberd API call.
type apiParams struct {
	name    string
	version int
	admin   bool // Does API require admin header ?

	method string
	query  url.Values
	body   []byte
}

//==============================================================================

type Stats struct {
	Name string `json:"name"`
}

type StatsResponse struct {
	Name string `json:"name"`
	Stat int    `json:"stat"`
}

func (StatsResponse) JSON() string {
	return "TODO"
}

func (s StatsResponse) String() string {
	return fmt.Sprintf("%d", s.Stat)
}

func (s *Stats) params() (apiParams, error) {
	var query url.Values
	if !stringInSlice(s.Name, s.knownStats()) {
		return apiParams{}, fmt.Errorf("unknow statistic: %s", s.Name)
	}

	body, err := json.Marshal(s)
	if err != nil {
		return apiParams{}, err
	}

	return apiParams{
		name:    "stats",
		version: 1,

		admin:  true,
		method: "POST",
		query:  query,
		body:   body,
	}, nil
}

func (s Stats) parseResponse(body []byte) (Response, error) {
	var resp StatsResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		// Cannot parse JSON response
		return ErrorResponse{Code: 99}, err
	}
	resp.Name = s.Name
	return resp, err
}

func (s Stats) knownStats() []string {
	return []string{"registeredusers", "onlineusers", "onlineusersnode",
		"uptimeseconds", "processes"}
}

//==============================================================================

type Register struct {
	JID      string `json:"jid"`
	Password string `json:"password"`
}

type RegisterResponse string

func (RegisterResponse) JSON() string {
	return "TODO"
}

func (r *Register) params() (apiParams, error) {
	var query url.Values

	jid, err := parseJID(r.JID)
	if err != nil {
		return apiParams{}, err
	}

	// Actual parameters for ejabberd. We expose JID string as it is
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
		return apiParams{}, err
	}

	return apiParams{
		name:    "register",
		version: 1,
		admin:   true,

		method: "POST",
		query:  query,
		body:   body,
	}, nil
}

func (r Register) parseResponse(body []byte) (Response, error) {
	var resp RegisterResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		// Cannot parse JSON response
		return ErrorResponse{Code: 99}, err
	}
	return resp, nil
}

//==============================================================================

type ErrorResponse struct {
	Code        int
	Description string
}

func (e ErrorResponse) JSON() string {
	return "TODO"
}

func (e ErrorResponse) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Description)
}
