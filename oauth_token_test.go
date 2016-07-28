package ejabberd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_GetToken(t *testing.T) {
	accessToken := "12345"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"access_token": "`+accessToken+`"}`)
	}))
	defer server.Close()

	// Make a transport that reroutes all traffic to the example server
	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	client := Client{URL: "http://localhost:5281", HTTPClient: &http.Client{Transport: transport}}
	token, err := client.GetToken("admin@localhost", "passw0rd", "sasl-auth", 3600)
	if err != nil {
		t.Errorf("GetToken failed: %s", err)
	}

	if token.AccessToken != accessToken {
		t.Errorf("Incorrect access token  %s != %s", token.AccessToken, accessToken)
	}
}
