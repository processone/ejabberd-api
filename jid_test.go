package ejabberd

import "testing"

func Test_JID(t *testing.T) {
	var tests = []struct {
		input string
		want  jid
		err   error
	}{
		{"username@domain/resource", jid{username: "username", domain: "domain", resource: "resource"}, nil},
		{"username@domain", jid{username: "username", domain: "domain"}, nil},
		{"username@domain/resourcewith/@", jid{username: "username", domain: "domain", resource: "resourcewith/@"}, nil},
		{"user@name@domain/resource", jid{}, testError{}},
	}
	for _, test := range tests {
		var got jid
		var err error
		if got, err = parseJID(test.input); err != nil && test.err == nil {
			t.Errorf("error on parseJID(%q): %s", test.input, err)
			return
		}
		if got != test.want {
			t.Errorf("parseJID(%q) = %q", test.input, got)
		}
	}
}
