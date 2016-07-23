package ejabberd

import (
	"errors"
	"fmt"
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
func GetToken(endpoint, sjid, password, scope, clientID string, duration time.Duration) (string, time.Time, error) {
	var j jid
	var t oauthToken
	var err error

	if j, err = parseJID(sjid); err != nil {
		return t.accessToken, t.expiration, err
	}

	var u string
	if u, err = JoinURL(endpoint, "authorization_token"); err != nil {
		return t.accessToken, t.expiration, err
	}

	ttl := int(duration.Seconds())

	now := time.Now()

	if t, err = httpGetToken(j, password, clientID, scope, strconv.Itoa(ttl), u); err != nil {
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

func httpGetToken(j jid, password, clientID, scope, ttl, apiURL string) (oauthToken, error) {
	params := params(j, password, clientID, scope, ttl)

	var errRedirectAttempt = errors.New("redirect")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return errRedirectAttempt
		},
	}

	var t oauthToken
	resp, err := client.PostForm(apiURL, params)

	// We expect a redirect on success: Check error for redirect:
	if urlError, ok := err.(*url.Error); ok && urlError.Err == errRedirectAttempt && resp.StatusCode == 302 {
		redirectURL := resp.Header.Get("Location")

		u, err := url.Parse(redirectURL)
		if err != nil {
			return t, err
		}

		result := url.Values{}
		if result, err = url.ParseQuery(u.RawQuery); err != nil {
			return t, err
		}

		if len(result["access_token"]) > 0 {
			t.accessToken = result["access_token"][0]
		}

		if len(result["error"]) > 0 {
			t.error = result["error"][0]
		}

		if len(result["expires_in"]) > 0 {
			expiresIn := result["expires_in"][0]
			var i int64
			if i, err = strconv.ParseInt(expiresIn, 10, 32); err != nil {
				return t, err
			}
			t.expiration = time.Now().Add(time.Duration(i) * time.Second)
		}

		resp.Body.Close()
		return t, nil
	}

	if err != nil {
		return t, fmt.Errorf("could not retrieve token: %s", err)
	}
	resp.Body.Close()

	if resp.StatusCode == 404 {
		return t, errors.New("oauth endpoint not found (404)")
	}

	return t, errors.New("unexpected reply from oauth endpoint")
}

func params(j jid, password, clientID, scope, ttl string) url.Values {
	return url.Values{
		"response_type": {"token"},
		"state":         {""},
		"client_id":     {clientID},
		"redirect_uri":  {""},
		"scope":         {scope},
		"username":      {j.bare()},
		"password":      {password},
		"ttl":           {ttl},
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
