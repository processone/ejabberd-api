package main

import (
	"fmt"
	"os"
	"time"

	"github.com/processone/ejabberd-api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app  = kingpin.New("ejabberd", "A command-line front-end for ejabberd server API.").Version("0.0.1").Author("ProcessOne")
	file = app.Flag("file", "OAuth token JSON file.").Short('f').Default(".ejabberd-oauth.json").String()

	// ========= token =========
	token         = app.Command("token", "Request an OAuth token.")
	tokenJID      = token.Flag("jid", "JID of the user to generate token for.").Short('j').Required().String()
	tokenPassword = token.Flag("password", "Password to use to retrieve user token.").Short('p').String()
	tokenAskPass  = token.Flag("prompt", "Prompt for password.").Short('P').Bool()
	tokenScope    = token.Flag("scope", "Comma separated list of scope to associate to token").Short('s').Default("sasl_auth").String()
	tokenTTL      = token.Flag("ttl", "Time before token expiration. Valid unit time are second (s), minutes (m), hours (h)").Default("8760h").Short('t').Duration()
	tokenEndpoint = token.Flag("endpoint", "ejabberd API endpoint.").Short('e').Default("http://localhost:5281/").String()
	tokenOauthURL = token.Flag("oauth-url", "Oauth suffix for oauth endpoint.").Default("/oauth/").String()

	// ========= stats =========
	stats     = app.Command("stats", "Get ejabberd statistics.")
	statsName = stats.Arg("name", "Name of stats to query.").Required().Enum("registeredusers", "onlineusers", "onlineusersnode", "uptimeseconds", "processes")

	// ========= user =========
	user          = app.Command("user", "Operations to perform on users.")
	userOperation = user.Arg("operation", "Operation").Required().Enum("register")
	userJID       = user.Flag("jid", "JID of the user to perform operation on.").Short('j').Required().String()
	userPassword  = user.Flag("password", "User password").Short('p').Required().String()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h') // BUG(mr) Short help flag does not seem to work.
	kingpin.CommandLine.Help = "A command-line front-end for ejabberd server API."

	command := kingpin.MustParse(app.Parse(os.Args[1:]))
	switch command {
	case token.FullCommand():
		getToken()
	default:
		execute(command)
	}
}

func execute(command string) {
	f, err := ejabberd.ReadOAuthFile(*file)
	if err != nil {
		kingpin.Fatalf("could not load token file %q: %s", *file, err)
	}

	c := ejabberd.Client{
		URL:   f.Endpoint + "api/",
		Token: f.AccessToken,
	}

	switch command {
	case stats.FullCommand():
		statsCommand(c)
	case user.FullCommand():
		userCommand(c, *userOperation)
	}
}

// =============================================================================

func getToken() {
	var token string
	var expiration time.Time
	var url string
	var err error
	if url, err = ejabberd.JoinURL(*tokenEndpoint, *tokenOauthURL); err != nil {
		kingpin.Fatalf("invalid endpoint URL: %s", err)
	}
	scope := ejabberd.PrepareScope(*tokenScope)
	if token, expiration, err = ejabberd.GetToken(url, *tokenJID, *tokenPassword, scope, *tokenTTL); err != nil {
		kingpin.Fatalf("could not retrieve token: %s", err)
	}

	var f ejabberd.OAuthFile
	f.AccessToken = token
	f.JID = *tokenJID
	f.Scope = scope
	f.Expiration = expiration
	f.Endpoint = *tokenEndpoint
	if err = f.Save(*file); err != nil {
		kingpin.Fatalf("could not save token to file %q: %s", *file, err)
	}
	fmt.Println("Successfully saved token in file", *file)
}

//==============================================================================

func statsCommand(c ejabberd.Client) {
	command := ejabberd.GetStats{
		Name: *statsName,
	}

	resp, err := c.Call(&command)
	if err != nil {
		kingpin.Fatalf("stats command error %q: %s", command.Name, err)
	}
	fmt.Println(string(resp))
}

//==============================================================================

func userCommand(c ejabberd.Client, op string) {
	switch op {
	case "register":
		registerCommand(c, *userJID, *userPassword)
	}
}

func registerCommand(c ejabberd.Client, j, p string) {
	// TODO Should we create a v2 command with only two parameters (JID, Password)
	command := ejabberd.RegisterUser{
		JID:      j,
		Password: p}
	resp, err := c.Call(&command)
	if err != nil {
		kingpin.Fatalf("register command error %v: %s", command, err)
	}
	fmt.Println(string(resp))
}

// TODO Interface for command result formatting
