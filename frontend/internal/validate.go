package frontend

import (
	"fmt"
	"regexp"
)

var (
	// Regular expression to validate username.
	regexpUsername          = regexp.MustCompile(`^[A-Za-z0-9-_]*$`)
	usernameValidCharacters = "latin letters, digits and characters '-' and '_'"
	minLengthUsername       = 3
	maxLengthUsername       = 32
	minLengthPassword       = 3
	maxLengthPassword       = 64
)

type ErrorFrontend struct {
	Err      error
	Frontend string
}

var (
	ErrorUsernameUnsupportedCharacters = ErrorFrontend{
		Err:      fmt.Errorf("username contains unsupported characters"),
		Frontend: fmt.Sprintf("Use %s for username.", usernameValidCharacters),
	}
	ErrorUsernameWrongLength = ErrorFrontend{
		Err:      fmt.Errorf("username wrong length"),
		Frontend: fmt.Sprintf("Username must be between %d and %d characters long.", minLengthUsername, maxLengthUsername),
	}
	ErrorPasswordWrongLength = ErrorFrontend{
		Err:      fmt.Errorf("password wrong length"),
		Frontend: fmt.Sprintf("Password must be between %d and %d characters long.", minLengthPassword, maxLengthPassword),
	}
)

func (e ErrorFrontend) Error() string {
	return e.Err.Error()
}

// ValidateUsername validates the username for valid length and for dangerous characters.
func ValidateUsername(u string) error {
	if !regexpUsername.MatchString(u) {
		return ErrorUsernameUnsupportedCharacters
	}
	if len(u) < minLengthUsername || len(u) > maxLengthUsername {
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
