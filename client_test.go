package ejabberd_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/processone/ejabberd-api"
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

	client := ejabberd.Client{BaseURL: "http://localhost:5281", HTTPClient: &http.Client{Transport: transport}}
	token, err := client.GetToken("admin@localhost", "passw0rd", "ejabberd:admin", 3600)
	if err != nil {
		t.Errorf("GetToken failed: %s", err)
	}

	if token.AccessToken != accessToken {
		t.Errorf("Incorrect access token  %s != %s", token.AccessToken, accessToken)
	}
}

// TODO provide const to specify token duration

func ExampleClient_GetToken() {
	client := ejabberd.Client{BaseURL: "http://localhost:5281"}

	if token, err := client.GetToken("admin@localhost", "passw0rd", "ejabberd:admin", 3600); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Retrieved access token:", token.AccessToken)
	}
}

func ExampleClient_Stats() {
	if flag.Lookup("test.v") == nil {
		fmt.Println("normal run")
	} else {
		fmt.Println("run under go test")
	}

	t := ejabberd.OAuthToken{AccessToken: "XjlJg0KF2wagT0A5dcYghePl8npsiEic"}
	client := ejabberd.Client{BaseURL: "http://localhost:5281", Token: t}

	if stats, err := client.Stats("registeredusers"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(stats.Name, stats.Value)
	}
	// Output:
	// 1
}
