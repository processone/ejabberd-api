/*
Package ejabberd expose ejabberd API in Go programming language.

You need to have ejabberd properly configured to expose API.

Here are needed step to get a working ejabberd configuration:

- Using ppolv/ejabberd
- Cherry-picked: https://github.com/processone/ejabberd/commit/ce0d1704c6cc167c8bc891587952f78c55f979ad

ejabberd configuration (see https://docs.ejabberd.im/admin/guide/oauth/):

- Added request handler for oauth:
      "/oauth": ejabberd_oauth
- Added configuration:
    commands_admin_access: configure
    commands:
      - add_commands: user
    oauth_expire: 32000000
    oauth_access: all
- Enabled mod_admin_extra module
- Added my user as admin

Generate a token for that user. Example from Erlang:
ejabberd_oauth:oauth_issue_token("ejabberd:user;ejabberd:admin").

 */
package ejabberd