# ejabberd API library and command-line tool

[![Build Status](https://semaphoreci.com/api/v1/processone/ejabberd-api/branches/master/shields_badge.svg)](https://semaphoreci.com/processone/ejabberd-api)
[![Coverage Status](https://coveralls.io/repos/github/processone/ejabberd-api/badge.svg?branch=master)](https://coveralls.io/github/processone/ejabberd-api?branch=master)

This tool is composed of two components:

- A command-line tool to interact with ejabberd through ReST API calls
  from the command-line, from any server type or desktop (Linux, OSX,
  Windows).
- An implementation of ejabberd API client library in Go. It can be
  used to interact with ejabberd from backend applications developed
  in Go programming language.

## Installation

Both the library and the command-line tool can be installed from
source with `go get` tool.

If you have go installed in your environment you can install
`ejabberd` command-line tool with:

```bash
go install -v github.com/processone/ejabberd-api/cmd/ejabberd@master
```
Once the get completes, you should find your new `ejabberd` (or `ejabberd.exe`) executable sitting inside `$GOPATH/bin/`.

## Configuring ejabberd OAuth API

Before being able to interact with ejabberd API, you need to configure
ejabberd with OAuth support enabled. This is is documented in
[ejabberd OAuth support](https://docs.ejabberd.im/developer/ejabberd-api/oauth/).

Here are example entries to check / change in your ejabberd
configuration file:

1. Add a [listener](https://docs.ejabberd.im/admin/configuration/listen/) for OAuth and ReST API:

   ```yaml
   listen:
     -
       # Using a separate port for oauth and API to make it easy to protect it
       # differently than BOSH and Websocket HTTP interface.
       port: 5281
       # oauth and API only listen on localhost interface for security reason
       # You can set ip to 0.0.0.0 to open it widely, but be careful!
       ip: 127.0.0.1
       module: ejabberd_http
       request_handlers:
         /oauth: ejabberd_oauth
         /api: mod_http_api
   ```

2. You can then configure the OAuth commands you want to expose. Check
   [`api_permissions`](https://docs.ejabberd.im/developer/ejabberd-api/permissions/)
   to make sure ACL for passing commands as admins are set properly:

   ```yaml
   api_permissions:
     "admin access":
       who:
         oauth:
           scope: "ejabberd:admin"
           access:
             allow:
               user: admin@localhost
       what:
         - "register"
         - "change_password"
   
   # Tokens are valid for a year as default:
   oauth_expire: 31536000
   oauth_access: all
   ```

3. Finally, make sure to enable the modules that provide the commands you
   want to use, for example:

   ```yaml
   modules:
     mod_admin_extra: {}
   ```

## ejabberd command-line tool

The `ejabberd` command-line tool allows you to interact with your ejabberd server
by using ejabberd's ReST API. It relies on OAuth tokens and scope to define the command
the user will be allowed to call.

### Usage

1. Generating an OAuth token:

   To use the `ejabberd` command-line tool, you first need to generate an OAuth
   token.

   It can be done, for example, with the following command:

   ```bash
   ejabberd token -j admin@localhost -p mypassword -s ejabberd:admin
   ```

   This will generate a `.ejabberd-oauth.json` file containing your
   credentials. Keep the file secret, as it will grant access to command
   available in the requested scope on your behalf.

2. Calling ejabberd API from the command-line, using your token file. For example:

   ```bash
   ejabberd stats registeredusers
   ```

3. You can also try to call any available command thanks to the generic `call` command. Do not forget `-a` parameter for commands that requires admin rights. For example:

   ```bash
   cat register.json
   {"user":"test1", "host":"localhost", "password":"somePassW0rd"}
   
   ejabberd call --name register -a --data-file=register.json
   ```

### Generating Bash/ZSH completion

You can generate Bash completion with following command:

```bash
./ejabberd --completion-script-bash
```

You can generate ZSH completion with following command:

```bash
./ejabberd --completion-script-zsh
```

To be able to use completion for Bash, you can type or add in your
`bash_profile` (or equivalent):

```bash
eval "$(ejabberd --completion-script-bash)"
```

For ZSH, you can use:

```bash
eval "$(ejabberd --completion-script-zsh)"
```

### Available commands

* **token**: Get OAuth token. This is needed before calling others commands.
* **stats**: Retrieve some stats from ejabberd.

To get a full list of commands and their options:

```bash
ejabberd --help-long
```

### OAuth Token file format

As a default, the token is stored in a file called
`./.ejabberd-oauth.json` when using the command `token` and read from
the same file when you use any other commands.

Option `-f file` will let you point to another file.

The file contains a JSON structure with the following fields:

* **access_token**: Actual token value.
* **endpoint**: Base URL.
* **jid**: JID for which user the token was generated.
* **scope**: OAuth scope for which the token was generated.
* **expiration**: Expiration date for the token.

For example:

```json
{"access_token":"AaQTb0PUZqeZhFKYoaTQBb4KKkCTAolE",
 "endpoint":"http://localhost:5281/",
 "jid":"admin@localhost",
 "scope":"ejabberd:admin",
 "expiration":"2017-07-23T13:53:08.326421575+02:00"}
```

## Development

### Local build

```bash
go build -o ejabberd cmd/ejabberd/main.go
```

### Running tests

You can run tests from repository clone with command:

```bash
go test -race -v ./.
```
