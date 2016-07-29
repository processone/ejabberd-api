package ejabberd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// OAuthToken defines how to store ejabberd OAuth token and
// attributes.
type OAuthToken struct {
	// Actual token value retrieved from server
	AccessToken string `json:"access_token"`
	Endpoint    string `json:"endpoint"`

	// Parameters associated with the token, stored for reference
	JID        string    `json:"jid"`
	Scope      string    `json:"scope"`
	Expiration time.Time `json:"expiration"`
}

// Save writes ejabberd OAuth structure to file.
func (t OAuthToken) Save(file string) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, b, 0640)
}

// ReadOAuthToken reads the content of JSon OAuth token file and
// return proper OAuthToken structure.
func ReadOAuthToken(file string) (OAuthToken, error) {
	var t OAuthToken
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(data, &t)
	return t, err
}

//==============================================================================

//==============================================================================
// HTTP

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

// tokenParams prepares HTTP form to retrieve token
func tokenParams(j jid, password, scope, ttl string) url.Values {
	return url.Values{
		"grant_type": {"password"},
		// TODO It would be nice to have ejabberd password grant_type support client_id:
		// "client_id":     {clientID},
		"scope":    {scope},
		"username": {j.bare()},
		"password": {password},
		"ttl":      {ttl},
	}
}

// ====
// Process ejabberd HTTP token response

func parseTokenError(body []byte) error {
	type jsonError struct {
		Error       string `json:"error"`
		Description string `json:"error_description"`
	}
	var e jsonError

	if err := json.Unmarshal(body, &e); err != nil {
		return errors.New("bad request")
	}
	return errors.New(e.Description)
}

func parseTokenResponse(body []byte) (OAuthToken, error) {
	type jsonResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		ExpiresIn   int    `json:"expires_in"`
	}
	var r jsonResp

	if err := json.Unmarshal(body, &r); err != nil {
		return OAuthToken{}, err
	}

	var t OAuthToken
	t.AccessToken = r.AccessToken
	t.Expiration = time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)

	return t, nil
}

//==============================================================================
// Helpers

// tokenURL Generate URL endpoint for retrieve a token using password
// grant type.
func tokenURL(baseURL string, oauthPath string) (string, error) {
	var path string
	var err error

	if oauthPath == "" {
		path, err = joinURL(baseURL, "oauth")
	} else {
		path, err = joinURL(baseURL, oauthPath)
	}

	if err != nil {
		return baseURL, err
	}

	return joinURL(path, "token")
}
