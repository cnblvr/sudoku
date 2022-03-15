package sudoku

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/cnblvr/sudoku/data"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"time"
)

// createAuthCookie creates the cookie and writes to the writer.
func (srv *Service) createAuthCookie(w http.ResponseWriter, a *data.Auth) error {
	value, err := srv.securecookie.Encode("auth", a)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     "auth",
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
	return nil
}

// deleteAuthCookie writes an empty cookie to be deleted from the cookie container on the client side.
func deleteAuthCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
}

// generatePasswordSalt generates the password hash salt. Different for each user.
func generatePasswordSalt() string {
	buf := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, buf)
	return base64.StdEncoding.EncodeToString(buf)
}

// hashPassword returns a hash of the password using the concatenation of the password, salt and pepper.
// The 'bcrypt' algorithm is used.
func (srv *Service) hashPassword(password string, salt string) (string, error) {
	saltBts, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return "", err
	}
	passwordBts := append(append([]byte(password), saltBts...), srv.passwordPepper...)
	hash, err := bcrypt.GenerateFromPassword(passwordBts, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

// verifyPassword verifies the password submitted by the user.
func (srv *Service) verifyPassword(password string, salt string, hash string) (bool, error) {
	saltBts, err := base64.StdEncoding.DecodeString(salt)
	if err != nil {
		return false, err
	}
	hashBts, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return false, err
	}
	passwordBts := append(append([]byte(password), saltBts...), srv.passwordPepper...)
	return bcrypt.CompareHashAndPassword(hashBts, passwordBts) == nil, nil
}
