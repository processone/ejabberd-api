package ejabberd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

// Client is an ejabberd client API wrapper. It is used to manage
// ejabberd client API interactions.
type Client struct {
	BaseURL string
	Token   string

	// Extra & Advanced features
	OAuthPath  string
	APIPath    string
	HTTPClient *http.Client
}

//==============================================================================

// Generic Call functions

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
		c.HTTPClient = defaultHTTPClient(15 * time.Second)
	}

	var url string
	if url, err = apiURL(c.BaseURL, c.OAuthPath, p.name); err != nil {
		return []byte{}, err
	}
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

//==============================================================================

// ==== Token ====

// TODO Get token from local file

// GetToken calls ejabberd API to get a token for a given scope, given
// valid jid and password.  We also assume that the user has the right
// to generate a token. In case of doubt you need to check ejabberd
// access option 'oauth_access'.
func (c Client) GetToken(sjid, password, scope string, duration time.Duration) (OAuthToken, error) {
	var j jid
	var t OAuthToken
	var err error

	// Set default values
	if c.HTTPClient == nil {
		c.HTTPClient = defaultHTTPClient(15 * time.Second)
	}

	if j, err = parseJID(sjid); err != nil {
		return t, err
	}

	var u string
	if u, err = tokenURL(c.BaseURL, c.OAuthPath); err != nil {
		return t, err
	}

	// Prepare token call parameters
	ttl := int(duration.Seconds())
	params := tokenParams(j, password, scope, strconv.Itoa(ttl))

	// Request token from server
	if t, err = httpGetToken(c.HTTPClient, u, params); err != nil {
		return t, err
	}
	return t, nil
}

//==============================================================================

// Stats allows to query ejabberd for generic statistics. Supported statistic names are:
//
//     registeredusers
//     onlineusers
//     onlineusersnode
//     uptimeseconds
//     processes
func (c Client) Stats(s Stats) (StatsResponse, error) {
	result, err := c.Call(s)
	if err != nil {
		return StatsResponse{}, err
	}
	resp := result.(StatsResponse)
	return resp, nil
}

//==============================================================================

// Pass default timeout
func defaultHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: timeout,
		},
	}
}
