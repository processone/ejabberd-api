# ejabberd API command-line tool

ejabberd command-line tool allow interacting ejabberd with ejabberd
ReST API. It relies on OAuth tokens and scope to define the command
the user will be allowed to call.

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

* token: Get OAuth token. This is needed before calling others commands.
