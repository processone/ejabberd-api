package ejabberd

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Request is the common interface to all ejabberd requests. It is
// passed to the ejabberd.Client Call methods to get parameters to
// make the call and parse responses from the server.
type Request interface {
	params() (apiParams, error)
	parseResponse([]byte) (Response, error)
}

// Response is the interface for all ejabberd API call results.
type Response interface {
	JSON() string
}

// apiParams gathers all values needed by the client to encode actual
// ejabberd API call. An ejabberd API commands should return apiParams
// struct when being issued params call.
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

func (s Stats) params() (apiParams, error) {
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

func (r Register) params() (apiParams, error) {
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
		return ErrorResponse{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	return resp, nil
}

//==============================================================================

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseError(body []byte) (ErrorResponse, error) {
	var resp ErrorResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return ErrorResponse{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	return resp, nil
}

func (e ErrorResponse) JSON() string {
	return "TODO"
}

func (e ErrorResponse) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
