/*
Package ejabberd exposes ejabberd API in Go programming language.

Installation

Both the library and the command-line tool can be installed from
source with go get tool.

If you have go installed in your environment you can install ejabberd
command-line tool with:

    go get -u github.com/processone/ejabberd-api

Configuring ejabberd OAuth API

Before being able to interact with ejabberd API, you need to configure
ejabberd with OAuth support enabled. This is is documented in ejabberd
OAuth support documentation
(https://docs.ejabberd.im/admin/guide/oauth/).

Here are example entries to check / change in your ejabberd
configuration file:

1. Add a listener for OAuth and ReST API:

   listen:
     -
       # Using a separate port for oauth and API to make it easy to protect it
       # differently than BOSH and Websocket HTTP interface.
       port: 5281
       # oauth and API only listen on localhost interface for security reason
       # You can set ip to "0.0.0.0" to open it widely, but be careful!
       ip: "127.0.0.1"
       module: ejabberd_http
       request_handlers:
         "/oauth": ejabberd_oauth
         "/api": mod_http_api

2. You can then configure the OAuth commands you want to expose. Check
`commands_admin_access` to make sure ACL for passing commands as admins
are set properly:

   commands_admin_access:
     - allow:
       - user: "admin@localhost"
   commands:
     - add_commands: [user, admin, open]
   # Tokens are valid for a year as default:
   oauth_expire: 31536000
   oauth_access: all

3. Finally, make sure the modules, you need to use the command from
are enabled, for example:

   modules:
     mod_admin_extra: {}

Using ejabberd client API

To use ejabberd Go client API, you need to have a token. If you do not
have one, you can get one programatically using 'Client.GetToken'
function.

From there, you can use any ejabberd client function to call any
ejabberd available command.

Check 'ejabberd.Client' documentation for details on the currently
available commands.

OAuth Token file format

As a default, the token is stored in a file called
'./.ejabberd-oauth.json' when using the command 'token' and read from
the same file when you use any other commands.

Option '-f file' will let you point to another file.

The file contains a JSON structure with the following fields:

    access_token   Actual token value.
    endpoint       Base URL.
    jid            JID for which user the token was generated.
    scope          OAuth scope for which the token was generated.
    expiration     Expiration date for the token.

For example:

    {"access_token":"AaQTb0PUZqeZhFKYoaTQBb4KKkCTAolE",
     "endpoint":"http://localhost:5281/",
     "jid":"admin@localhost",
     "scope":"ejabberd:admin",
     "expiration":"2017-07-23T13:53:08.326421575+02:00"}

*/
package ejabberd
