package ejabberd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_GetOauthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"access_token": "12345"}`)
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
		t.Errorf("GetOAuthToken failed: %s", err)
	}

	fmt.Println(token)
}
