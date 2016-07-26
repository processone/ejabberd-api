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

type oauthToken struct {
	accessToken string
	expiration  time.Time
	error       string
}

// GetToken calls ejabberd API to get a for a given scope, given valid jid and password.
// We also assume that the user has the right to generate a token.
func GetToken(endpoint, sjid, password, scope string, duration time.Duration) (string, time.Time, error) {
	var j jid
	var t oauthToken
	var err error

	if j, err = parseJID(sjid); err != nil {
		return t.accessToken, t.expiration, err
	}

	var u string
	if u, err = JoinURL(endpoint, "token"); err != nil {
		return t.accessToken, t.expiration, err
	}

	ttl := int(duration.Seconds())

	now := time.Now()

	if t, err = httpGetToken(j, password, scope, strconv.Itoa(ttl), u); err != nil {
		return t.accessToken, t.expiration, err
	}

	if t.expiration.IsZero() {
		t.expiration = now.Add(duration)
	}

	if t.error != "" {
		return t.accessToken, t.expiration, fmt.Errorf(t.error)
	}

	return t.accessToken, t.expiration, nil
}

type jsonResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
}

type jsonError struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

func httpGetToken(j jid, password, scope, ttl, apiURL string) (oauthToken, error) {
	var t oauthToken

	client := &http.Client{}
	params := params(j, password, scope, ttl)
	resp, err := client.PostForm(apiURL, params)

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

	t.accessToken = r.AccessToken
	t.expiration = time.Now().Add(time.Duration(r.ExpiresIn) * time.Second)

	return t, nil
}

func params(j jid, password, scope, ttl string) url.Values {
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

// =============================================================================
