package main

import (
	"fmt"
	"io/ioutil"
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
	statsName = stats.Arg("name", "Name of stats to query.").Required().String()

	// ========= admin =========
	register         = app.Command("register", "Create a new user.")
	registerJID      = register.Flag("jid", "JID of the user to create.").Short('j').Required().String()
	registerPassword = register.Flag("password", "Password to set for created user.").Short('p').Required().String()

	// ========= user =========
	user          = app.Command("user", "Operations to perform on users.")
	userOperation = user.Arg("operation", "Operation").Required().Enum("resources")
	userJID       = user.Flag("jid", "JID of the user to perform operation on.").Short('j').String()

	// ========= offline =========
	offline          = app.Command("offline", "Operations to perform on offline store.")
	offlineOperation = offline.Arg("operation", "Operation").Required().Enum("count")
	offlineJID       = offline.Flag("jid", "JID of the user to perform operation on, if different from token owner").Short('j').String()

	// ========= generic call =========
	call      = app.Command("command", "Call a command on ejabberd server, using your token credentials.")
	callFile  = call.Flag("data-file", "File with JSON data to send to ejabberd. You can also use /dev/stdin").String()
	callData  = call.Flag("data", "File with JSON data to send to ejabberd. Omit to read from STDIN").String()
	callName  = call.Flag("name", "Name of command on server").Short('n').Required().String()
	callAdmin = call.Flag("admin", "Call as admin").Short('a').Bool()
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
		Token:   t,
	}

	switch command {
	case call.FullCommand():
		genericCommand(c, *callName, *callData, *callFile, *callAdmin)
	case register.FullCommand():
		registerCommand(c, *registerJID, *registerPassword)
	case stats.FullCommand():
		statsCommand(c)
	case user.FullCommand():
		userCommand(c, *userOperation)
	case offline.FullCommand():
		offlineCommand(c, *offlineOperation)
	}

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

func registerCommand(c ejabberd.Client, j, p string) {
	resp, err := c.RegisterUser(j, p)
	if err != nil {
		kingpin.Fatalf("user registration error for %s: %s", j, err)
	}
	format(resp)
}

//==============================================================================

func statsCommand(c ejabberd.Client) {
	resp, err := c.Stats(*statsName)
	if err != nil {
		kingpin.Fatalf("stats error %q: %s", *statsName, err)
	}
	format(resp)
}

//==============================================================================

func userCommand(c ejabberd.Client, op string) {
	switch op {
	case "resources":
		resourcesCommand(c, *userJID)
	}
}

func resourcesCommand(c ejabberd.Client, jid string) {
	if jid == "" {
		jid = c.Token.JID
	}

	resp, err := c.UserResources(jid)
	if err != nil {
		kingpin.Fatalf("%s: %s", jid, err)
	}
	format(resp)
}

//==============================================================================

func offlineCommand(c ejabberd.Client, op string) {
	switch op {
	case "count":
		offlineCountCommand(c, *offlineJID)
	}
}

func offlineCountCommand(c ejabberd.Client, jid string) {
	if jid == "" {
		jid = c.Token.JID
	}
	resp, err := c.GetOfflineCount(jid)
	if err != nil {
		kingpin.Fatalf("offline count error for %s: %s", jid, err)
	}
	format(resp)
}

//==============================================================================

func genericCommand(c ejabberd.Client, commandName, input string, file string, admin bool) {
	var data []byte
	var err error

	if file != "" && input != "" {
		kingpin.Fatalf("Use either data or data-file option to pass input to ejabberd API.")
	}

	if input != "" {
		data = []byte(input)
	} else {
		switch file {
		case "/dev/stdin":
			data, err = ioutil.ReadAll(os.Stdin)
		case "":
			// Some valid ejabberd commands accept empty input
		default:
			data, err = ioutil.ReadFile(file)
		}
	}

	if err != nil {
		kingpin.Fatalf("%s", err)
	}

	code, result, err := c.CallRaw(data, commandName, admin)
	if err != nil {
		kingpin.Fatalf("%s", err)
	}
	if code != 200 {
		fmt.Printf("Response: %d\n", code)
		kingpin.Fatalf("%s", result)
	}
	fmt.Printf("%s", result)
}
