package ejabberd

import "testing"

type testError struct{}

func (e testError) Error() string {
	return "expected error"
}

func TestJoinURL(t *testing.T) {
	var tests = []struct {
		baseURL string
		suffix  string
		want    string
		err     error
	}{
		{"localhost:5281", "", "", testError{}},
	}
	for _, test := range tests {
		var got string
		var err error
		if got, err = joinURL(test.baseURL, test.suffix); err != nil && test.err == nil {
			t.Errorf("error on JoinURL(%q, %q): %s", test.baseURL, test.suffix, err)
			return
		}
		if got != test.want {
			t.Errorf("JoinURL(%q, %q) = %q", test.baseURL, test.suffix, got)
		}
	}
}
