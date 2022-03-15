package data

import (
	"fmt"
	"regexp"
)

// Auth is the authorization information contained in user's cookie and is used to control handlers and templates.
type Auth struct {
	// Flag indicates if the user is logged in.
	IsAuthorized bool `json:"isAuthorized"`
	// Username in lower case.
	Username string `json:"username"`
}

var (
	// Regular expression to validate username.
	regexpUsername = regexp.MustCompile(`^[A-Za-z0-9-_]*$`)
)

var (
	ErrorUsernameUnsupportedCharacters = fmt.Errorf("username contains unsupported characters")
	ErrorUsernameWrongLength           = fmt.Errorf("username wrong length")
	ErrorPasswordWrongLength           = fmt.Errorf("password wrong length")
)

// ValidateUsername validates the username for valid length and for dangerous characters.
func ValidateUsername(u string) error {
	if !regexpUsername.MatchString(u) {
		return ErrorUsernameUnsupportedCharacters
	}
	if len(u) < 3 || len(u) > 32 {
		return ErrorUsernameWrongLength
	}
	return nil
}

// ValidatePassword validates the password for valid length.
func ValidatePassword(p string) error {
	if len(p) < 3 || len(p) > 64 {
		return ErrorPasswordWrongLength
	}
	return nil
}
