package ejabberd

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

// OAuthFile defines ejabberd OAuth file structure.
type OAuthFile struct {
	// Actual token value retrieved from server
	AccessToken string
	Endpoint    string

	// Parameters associated with the token, stored for reference
	JID        string
	Scope      string
	Expiration time.Time
}

// Save write ejabberd OAuth structure to file.
func (f OAuthFile) Save(file string) error {
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, b, 0640)
}

// ReadOAuthFile reads the content of JSon Oauth token file and return
// proper OAuthFile structure.
func ReadOAuthFile(file string) (OAuthFile, error) {
	var f OAuthFile
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return f, err
	}

	err = json.Unmarshal(data, &f)
	return f, err
}
