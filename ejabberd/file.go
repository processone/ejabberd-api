package ejabberd

import (
	"encoding/json"
	"io/ioutil"
)

// OAuthFile defines ejabberd OAuth file structure.
type OAuthFile struct {
	AccessToken string
	// Reminder of parameters associated with the token
	JID   string
	Scope string
}

// Save write ejabberd OAuth structure to file.
func (f OAuthFile) Save(file string) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, b, 0640)
}
