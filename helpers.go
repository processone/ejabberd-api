package ejabberd

import "strings"

// prepareScope ensures we return scopes as space separated. However,
// we accept comma separated scopes as input as well for convenience.
func prepareScope(s string) string {
	return strings.Replace(s, ",", " ", -1)
}

//==============================================================================
// Internal helper functions

// stringInSlice returns whether a string is a member of a string
// slice.
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
