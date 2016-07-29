package ejabberd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Client is an ejabberd client API wrapper. It is used to manage
// ejabberd client API interactions.
type Client struct {
	BaseURL    string
	OAuthPath  string
	HTTPClient *http.Client

	// TODO refactor
	Token string
}

// TODO Get token from local file

// GetToken calls ejabberd API to get a token for a given scope, given
// valid jid and password.  We also assume that the user has the right
// to generate a token. In case of doubt you need to check ejabberd
// access option 'oauth_access'.
func (c *Client) GetToken(sjid, password, scope string, duration time.Duration) (OAuthToken, error) {
	var j jid
	var t OAuthToken
	var err error

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
	}

	if j, err = parseJID(sjid); err != nil {
		return t, err
	}

	var u string
	if u, err = c.tokenURL(); err != nil {
		return t, err
	}

	ttl := int(duration.Seconds())
	now := time.Now()

	params := tokenParams(j, password, scope, strconv.Itoa(ttl))
	if t, err = httpGetToken(c.HTTPClient, u, params); err != nil {
		return t, err
	}

	if t.Expiration.IsZero() {
		t.Expiration = now.Add(duration)
	}

	return t, nil
}

//===============================

func httpGetToken(c *http.Client, apiURL string, params url.Values) (OAuthToken, error) {
	// Performs HTTP request
	resp, err := c.PostForm(apiURL, params)
	if err != nil {
		return OAuthToken{}, err
	}
	defer resp.Body.Close()

	// Endpoint not found
	if resp.StatusCode == 404 {
		return OAuthToken{}, errors.New("oauth endpoint not found (404)")
	}

	// Cannot read HTTP response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return OAuthToken{}, errors.New("cannot read HTTP response from server")
	}

	// Bad request
	if resp.StatusCode == 400 {
		return OAuthToken{}, parseTokenError(body)
	}

	// Success
	return parseTokenResponse(body)
}

//==============================================================================
// Helpers

func (c *Client) tokenURL() (string, error) {
	var path string
	var err error

	if c.OAuthPath == "" {
		path, err = joinURL(c.BaseURL, "oauth")
	} else {
		path, err = joinURL(c.BaseURL, c.OAuthPath)
	}

	if err != nil {
		return c.BaseURL, err
	}

	return joinURL(path, "token")
}
