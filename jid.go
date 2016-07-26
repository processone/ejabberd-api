package ejabberd

import (
	"errors"
	"fmt"
	"strings"
)

// JID processing
// TODO update gox and import it directly from gox

type jid struct {
	username string
	domain   string
	resource string
}

func parseJID(sjid string) (jid, error) {
	var j jid

	s1 := strings.SplitN(sjid, "/", 2)
	if len(s1) > 1 {
		j.resource = s1[1]
	}

	s2 := strings.Split(s1[0], "@")
	if len(s2) != 2 {
		return jid{}, errors.New("invalid jid")
	}

	j.username = s2[0]
	j.domain = s2[1]
	return j, nil
}

func (j jid) bare() string {
	return fmt.Sprintf("%s@%s", j.username, j.domain)
}
