package ejabberd

import (
	"net/http"
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

	// Set default values
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{}
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
