package main

import (
	"fmt"
	"os"

	"github.com/processone/ejabberd-api"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app  = kingpin.New("ejabberd", "A command-line front-end for ejabberd server API.").Version("0.0.1").Author("ProcessOne")
	file = app.Flag("file", "OAuth token JSON file.").Short('f').Default(".ejabberd-oauth.json").String()
	json = app.Flag("json", "JSON formatted output").Bool()

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
	t, err := ejabberd.ReadOAuthToken(*file)
	if err != nil {
		kingpin.Fatalf("could not load token file %q: %s", *file, err)
	}
	if t.AccessToken == "" {
		kingpin.Fatalf("could not find access_token in file %q", *file)
	}
	c := ejabberd.Client{
		BaseURL: t.Endpoint,
		APIPath: "api/",
		Token:   t.AccessToken,
	}

	var resp ejabberd.Response
	switch command {
	case stats.FullCommand():
		resp = statsCommand(c)
	case user.FullCommand():
		resp = userCommand(c, *userOperation)
	}
	format(resp)
}

func format(resp ejabberd.Response) {
	if *json {
		fmt.Println(resp.JSON())
	} else {
		fmt.Println(resp)
	}
}

// =============================================================================

func getToken() {
	var token ejabberd.OAuthToken
	var err error
	client := ejabberd.Client{BaseURL: *tokenEndpoint, OAuthPath: *tokenOauthURL}
	if token, err = client.GetToken(*tokenJID, *tokenPassword, *tokenScope, *tokenTTL); err != nil {
		kingpin.Fatalf("could not retrieve token: %s", err)
	}

	token.JID = *tokenJID
	token.Endpoint = *tokenEndpoint
	if err = token.Save(*file); err != nil {
		kingpin.Fatalf("could not save token to file %q: %s", *file, err)
	}
	fmt.Println("Successfully saved token in file", *file)
}

//==============================================================================

func statsCommand(c ejabberd.Client) ejabberd.Response {
	command := ejabberd.Stats{
		Name: *statsName,
	}

	resp, err := c.Stats(command)
	if err != nil {
		kingpin.Fatalf("stats command error %q: %s", command.Name, err)
	}
	return resp
}

//==============================================================================

func userCommand(c ejabberd.Client, op string) ejabberd.Response {
	var resp ejabberd.Response
	switch op {
	case "register":
		resp = registerCommand(c, *userJID, *userPassword)
	}
	return resp
}

func registerCommand(c ejabberd.Client, j, p string) ejabberd.Response {
	// TODO Should we create a v2 command with only two parameters (JID, Password)
	command := ejabberd.Register{
		JID:      j,
		Password: p}
	resp, err := c.Call(&command)
	if err != nil {
		kingpin.Fatalf("register command error %v: %s", command, err)
	}
	return resp
}

// TODO Interface for command result formatting
