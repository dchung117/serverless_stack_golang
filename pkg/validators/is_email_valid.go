package validators

import (
	"regexp"
)

func IsEmailValid(email string) bool {
	// parse out regex email from email string
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]{1,64}@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// determine if email is invalid
	if len(email) < 3 || len(email) > 254 || !rxEmail.MatchString(email) {
		return false
	}

	return true
}
