package ejabberd

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

//==============================================================================
// Helpers for command-line tool

// joinURL checks that Base URL is a valid URL and joins base URL with
// the method suffix string.
func joinURL(baseURL string, suffix string) (string, error) {
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

//==============================================================================
// Internal helper functions

// stringInSlice returns whether a string is a member of a string
// slice.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
