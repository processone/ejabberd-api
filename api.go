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

// TODO: Move into a api_stats file

// Wraps various ejabberd call that all returns stats
// From ejabberd mod_admin_extra

type Stats struct {
	Name string `json:"name"`
	Host string `json:"host,omitempty"`
}

type StatsResponse struct {
	Name  string `json:"name"`
	Host  string `json:"host,omitempty"`
	Value int    `json:"stat"`
}

func (s StatsResponse) JSON() string {
	body, _ := json.Marshal(s)
	return string(body)
}

func (s StatsResponse) String() string {
	return fmt.Sprintf("%d", s.Value)
}

func (s Stats) params() (apiParams, error) {
	switch s.Name {
	case "":
		return apiParams{}, fmt.Errorf("required argument 'name' not provided")
	case "registeredusers", "onlineusers", "onlineusersnode", "uptimeseconds", "processes":
		return s.paramsStats()
	default:
		return apiParams{}, fmt.Errorf("unknow statistic: %s", s.Name)
	}
}

func (s Stats) paramsStats() (apiParams, error) {
	var query url.Values

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
		return ErrorResponse{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	resp.Name = s.Name
	return resp, err
}

//==============================================================================

// From ejabberd_admin

type Register struct {
	JID      string `json:"jid"`
	Password string `json:"password"`
}

type RegisterResponse string

func (r RegisterResponse) JSON() string {
	body, _ := json.Marshal(r)
	return string(body)
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

type OfflineCount struct {
	JID string `json:"jid"`
}

type OfflineCountResponse struct {
	Name  string `json:"name"`
	JID   string `json:"jid"`
	Value int    `json:"value"`
}

func (o OfflineCountResponse) JSON() string {
	body, _ := json.Marshal(o)
	return string(body)
}

func (o OfflineCount) params() (apiParams, error) {
	var query url.Values
	jid, err := parseJID(o.JID)
	if err != nil {
		return apiParams{}, err
	}

	type offlineCount struct {
		User   string `json:"user"`
		Server string `json:"server"`
	}

	data := offlineCount{
		User:   jid.username,
		Server: jid.domain,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return apiParams{}, err
	}

	if err != nil {
		return apiParams{}, err
	}

	return apiParams{
		name:    "get_offline_count",
		version: 1,

		method: "POST",
		query:  query,
		body:   body,
	}, nil
}

func (o OfflineCount) parseResponse(body []byte) (Response, error) {
	var resp OfflineCountResponse
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return ErrorResponse{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	resp.Name = "offline_count"
	resp.JID = o.JID
	return resp, nil
}

func (o OfflineCountResponse) String() string {
	return fmt.Sprintf("%d", o.Value)
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
	body, _ := json.Marshal(e)
	return string(body)
}

func (e ErrorResponse) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
