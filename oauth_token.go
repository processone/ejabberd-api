package ejabberd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
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

// Helpers for command-line tool

// JoinURL checks that Base URL is a valid URL and joins base URL with
// the method suffix string.
func JoinURL(baseURL string, suffix string) (string, error) {
	var u *url.URL
	var err error

	if u, err = url.Parse(baseURL); err != nil {
		return "", fmt.Errorf("invalid url: %s", baseURL)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("invalid url scheme: %s", u.Scheme)
	}

	u.Path = path.Join(u.Path, suffix)
	return u.String(), nil
}

// PrepareScope ensures we return scopes as space separated. However,
// we accept comma separated scopes as input as well for convenience.
func PrepareScope(s string) string {
	return strings.Replace(s, ",", " ", -1)
}

// =============================================================================

// JID processing
// TODO update gox and import it directly from gox

type jid struct {
	username string
	domain   string
	resource string
}

func parseJID(sjid string) (jid, error) {
	var j jid

	s1 := strings.SplitN(sjid, "/", 2)
	if len(s1) > 1 {
		j.resource = s1[1]
	}

	s2 := strings.Split(s1[0], "@")
	if len(s2) != 2 {
		return jid{}, errors.New("invalid jid")
	}

	j.username = s2[0]
	j.domain = s2[1]
	return j, nil
}

func (j jid) bare() string {
	return fmt.Sprintf("%s@%s", j.username, j.domain)
}
