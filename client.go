package ejabberd

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (c *Client) tokenURL() (string, error) {
	var path string
	var err error

	if c.OAuthPath == "" {
		path, err = JoinURL(c.BaseURL, "oauth")
	} else {
		path, err = JoinURL(c.BaseURL, c.OAuthPath)
	}

	if err != nil {
		return c.BaseURL, err
	}

	return JoinURL(path, "token")
}

// TODO Get token from local file

// GetToken calls ejabberd API to get a token for a given scope, given
// valid jid and password.  We also assume that the user has the right
// to generate a token. In case of doubt you need to check ejabberd
// access option `oauth_access`.
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

	params := params(j, password, scope, strconv.Itoa(ttl))
	if t, err = httpGetToken(c.HTTPClient, u, params); err != nil {
		return t, err
	}

	if t.Expiration.IsZero() {
		t.Expiration = now.Add(duration)
	}

	if t.error != "" {
		return t, fmt.Errorf(t.error)
	}

	return t, nil
}

//===============================
// http_api.go

func httpGetToken(c *http.Client, apiURL string, params url.Values) (OAuthToken, error) {
	var t OAuthToken

	resp, err := c.PostForm(apiURL, params)

	if err != nil {
		return t, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return t, errors.New("oauth endpoint not found (404)")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 400 {
		var e jsonError
		if err := json.Unmarshal(body, &e); err != nil {
			return t, errors.New("bad request")
		}
		return t, errors.New(e.Description)
	}

	var r jsonResp
	if err := json.Unmarshal(body, &r); err != nil {
		return t, err
	}

	t.AccessToken = r.AccessToken
	t.Expiration = time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)

	return t, nil
}
