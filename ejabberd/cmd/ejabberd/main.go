package main

import (
	"fmt"
	"os"

	"github.com/processone/go/ejabberd"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app      = kingpin.New("ejabberd", "A command-line front-end for ejabberd server API.").Version("0.0.1").Author("ProcessOne")
	endpoint = app.Flag("endpoint", "ejabberd API endpoint.").Short('e').Default("http://localhost:5281/").String()
	file     = app.Flag("file", "OAuth token JSON file.").Short('f').Default(".ejabberd-oauth.json").String()

	// ========= token =========
	token         = app.Command("token", "Request an OAuth token.")
	tokenJID      = token.Flag("jid", "JID of the user to generate token for.").Short('j').Required().String()
	tokenPassword = token.Flag("password", "Password to use to retrieve user token.").Short('p').String()
	tokenAskPass  = token.Flag("prompt", "Prompt for password.").Short('P').Bool()
	tokenScope    = token.Flag("scope", "Comma separated list of scope to associate to token").Short('s').String()
	tokenClient   = token.Flag("client", "Name of application that will use the token.").Default("go-ejabberd").String()
	tokenOauthURL = token.Flag("oauth-url", "Oauth suffix for oauth endpoint.").Default("/oauth/").String()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h') // BUG(mr) Short help flag does not seem to work.
	kingpin.CommandLine.Help = "A command-line front-end for ejabberd server API."

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case token.FullCommand():
		getToken()
	}
}

// =============================================================================

func getToken() {
	var token string
	var url string
	var err error
	if url, err = ejabberd.JoinURL(*endpoint, *tokenOauthURL); err != nil {
		kingpin.Fatalf("invalid endpoint URL: %s", err)
	}
	if token, err = ejabberd.GetToken(url, *tokenJID, *tokenPassword, "get_roster sasl_auth", *tokenClient); err != nil {
		kingpin.Fatalf("could not retrieve token: %s", err)
	}

	var f ejabberd.OAuthFile
	f.AccessToken = token
	f.JID = *tokenJID
	if err = f.Save(*file); err != nil {
		kingpin.Fatalf("could not save token to file %q: %s", *file, err)
	}
	fmt.Println("Successfully saved token in file", *file)
}
