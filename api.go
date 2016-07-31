package ejabberd

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Response is the common interface for all ejabberd API call results.
type Response interface {
	JSON() string
}

// request is the common interface to all ejabberd requests. It is
// passed to the ejabberd.Client Call methods to get parameters to
// make the call and parse responses from the server.
type request interface {
	params() (apiParams, error)
	parseResponse([]byte) (Response, error)
}

// apiParams gathers all values needed by the client to encode actual
// ejabberd API call. An ejabberd API commands should return apiParams
// struct when being issued params call.
type apiParams struct {
	name    string
	version int
	admin   bool // = Flag to mark if API requires admin header

	method string
	query  url.Values
	body   []byte
}

//==============================================================================

// TODO: Move into a api_stats file

// Wraps various ejabberd call that all returns stats
// From ejabberd mod_admin_extra

// Stats is the data structure returned by ejabberd Stats API call.
type Stats struct {
	Name  string `json:"name"`
	Value int    `json:"stat"`
}

// JSON converts Stats data structure to JSON string.
func (s Stats) JSON() string {
	body, _ := json.Marshal(s)
	return string(body)
}

// String represents Stats data structure as a human readable value.
func (s Stats) String() string {
	return fmt.Sprintf("%d", s.Value)
}

type statsRequest struct {
	Name string `json:"name"`
}

func (s statsRequest) params() (apiParams, error) {
	switch s.Name {
	case "":
		return apiParams{}, fmt.Errorf("required argument 'name' not provided")
	case "registeredusers", "onlineusers", "onlineusersnode", "uptimeseconds", "processes":
		return s.paramsStats()
	default:
		return apiParams{}, fmt.Errorf("unknow statistic: %s", s.Name)
	}
}

func (s statsRequest) paramsStats() (apiParams, error) {
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

func (s statsRequest) parseResponse(body []byte) (Response, error) {
	var resp Stats
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return Error{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	resp.Name = s.Name
	return resp, err
}

//==============================================================================

// From ejabberd_admin

// Register contains the message return by server after successful
// user registration.
type Register string

// JSON represents Register result as a JSON string. Can be useful for
// further processing.
func (r Register) JSON() string {
	body, _ := json.Marshal(r)
	return string(body)
}

type registerRequest struct {
	JID      string `json:"jid"`
	Password string `json:"password"`
}

func (r registerRequest) params() (apiParams, error) {
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

func (r registerRequest) parseResponse(body []byte) (Response, error) {
	var resp Register
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return Error{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	return resp, nil
}

//==============================================================================

// OfflineCount contains the result of the call to ejabberd
// get_offline_count API.
type OfflineCount struct {
	Name  string `json:"name"`
	JID   string `json:"jid"`
	Value int    `json:"value"`
}

// JSON represents OfflineCount as a JSON string, for further
// processing with other tools.
func (o OfflineCount) JSON() string {
	body, _ := json.Marshal(o)
	return string(body)
}

type offlineCountRequest struct {
	JID string `json:"jid"`
}

func (o offlineCountRequest) params() (apiParams, error) {
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

func (o offlineCountRequest) parseResponse(body []byte) (Response, error) {
	var resp OfflineCount
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return Error{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	resp.Name = "offline_count"
	resp.JID = o.JID
	return resp, nil
}

func (o OfflineCount) String() string {
	return fmt.Sprintf("%d", o.Value)
}

//==============================================================================

// Error represents ejabberd error returned by the server as result of
// ejabberd API calls.
type Error struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func parseError(body []byte) (Error, error) {
	var resp Error
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return Error{Code: 99, Message: "Cannot parse JSON response"}, err
	}
	return resp, nil
}

// JSON represents ejabberd error response as a JSON string, for further
// processing with other tools.
func (e Error) JSON() string {
	body, _ := json.Marshal(e)
	return string(body)
}

func (e Error) String() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
