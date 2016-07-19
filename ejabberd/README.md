# ejabberd API command-line tool

ejabberd command-line tool allow interacting ejabberd with ejabberd
ReST API. It relies on OAuth tokens and scope to define the command
the user will be allowed to call.

## Configure ejabberd OAuth API

You need to configure ejabberd with OAuth support enabled. This is is
documented in
[ejabberd OAuth support](https://docs.ejabberd.im/admin/guide/oauth/).

Here are example point to check / change in configuration:

1. Add a listener for OAuth and ReST API:

   ```yaml
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
   ```

2. You can then configure the OAuth commands you want to expose. Check
   `commands_admin_access` to make sure ACL for passing commands as
   admins are set properly:

   ```yaml
   commands_admin_access:
     - allow:
       - user: "admin@localhost"
   commands:
     - add_commands: ['ejabberd:user', user, admin, open]
   oauth_expire: 32000000
   oauth_access: all
   ```

3. Finally, make sure the modules, you need to use the command from
   are enabled, for example:

   ```yaml
   modules:
     mod_admin_extra: {}
   ```

## Installation

The tool can be installed from source with `go get` tool.

If you have go installed in your environment you can install
`ejabberd` command-line tool with:

```bash
go get -u github.com/processone/go/ejabberd/cmd/ejabberd
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

## Commands

* **token**: Get OAuth token. This is needed before calling others commands.
