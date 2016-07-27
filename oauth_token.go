package ejabberd

import (
	"net/url"
	"time"
)

type OAuthToken struct {
	AccessToken string
	Expiration  time.Time
	error       string
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
