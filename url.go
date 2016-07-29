package ejabberd

import (
	"fmt"
	"net/url"
	"path"
)

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

// tokenURL generates URL endpoint for retrieving a token using
// password grant type.
func tokenURL(baseURL, oauthPath string) (string, error) {
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

// apiURL generates URL endpoint for calling a given ejabberd API
// command name.
func apiURL(baseURL, apiPath, name string) (string, error) {
	var path string
	var err error

	if apiPath == "" {
		path, err = joinURL(baseURL, "api")
	} else {
		path, err = joinURL(baseURL, apiPath)
	}

	if err != nil {
		return baseURL, err
	}

	return joinURL(path, name+"/")
}
