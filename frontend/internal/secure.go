package frontend

import (
	"crypto/rand"
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
		Expires:  time.Now().Add(time.Hour * 24),
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
func generatePasswordSalt() []byte {
	buf := make([]byte, 16)
	_, _ = io.ReadFull(rand.Reader, buf)
	return buf
}

// hashPassword returns a hash of the password using the concatenation of the password, salt and pepper.
// The 'bcrypt' algorithm is used.
func (srv *Service) hashPassword(password string, salt []byte) ([]byte, error) {
	passwordBts := append(append([]byte(password), salt...), srv.passwordPepper...)
	hash, err := bcrypt.GenerateFromPassword(passwordBts, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// verifyPassword verifies the password submitted by the user.
func (srv *Service) verifyPassword(password string, salt, hash []byte) (bool, error) {
	passwordBts := append(append([]byte(password), salt...), srv.passwordPepper...)
	return bcrypt.CompareHashAndPassword(hash, passwordBts) == nil, nil
}
