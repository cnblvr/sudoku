package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/model"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

const (
	ErrorBadRequest                 = "Bad request."
	ErrorInternalServerError        = "Internal server error."
	ErrorUsernameOrPasswordNotValid = "Username or password not valid."
	ErrorUsernameAlreadyTaken       = "Username already taken."
	ErrorUsernameIsNotValid         = "Username is not valid."
	ErrorPasswordIsNotValid         = "Password is not valid."
	ErrorUsernamePasswordSame       = "Username must not be the same as password."
	ErrorPasswordsMustMatch         = "Passwords must match."
)

// HandleLogin is a handler of login page.
func (srv *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	redirect := func(endpoint string) {
		http.Redirect(w, r, endpoint, http.StatusSeeOther)
	}

	// If the user is already logged in, then redirect to the main page.
	auth := getAuth(r)
	if auth.IsAuthorized {
		log.Debug().Str("redirect", data.EndpointIndex).Int64("id", auth.ID).Msg("client already logged in")
		redirect(data.EndpointIndex)
		return
	}

	var Data struct {
		ErrorMessage string
	}

	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("login"),
		},
	}

	// POST method processes data from the user
	if r.Method == http.MethodPost {
		Data.ErrorMessage = func() string {
			if err := r.ParseForm(); err != nil {
				log.Warn().Err(err).Msg("failed to parse form")
				return ErrorBadRequest
			}
			username, password := r.Form.Get("_username"), r.Form.Get("_password")
			if err := ValidateUsername(username); err != nil {
				return ErrorUsernameOrPasswordNotValid
			}
			user, isExists, err := model.UserByUsername(srv.redis, username)
			if err != nil {
				log.Error().Err(err).Msg("failed to get user")
				return ErrorInternalServerError
			}
			if !isExists {
				log.Debug().Err(err).Msg("username is not exists")
				return ErrorUsernameOrPasswordNotValid
			}
			auth = &data.Auth{
				IsAuthorized: true,
				ID:           user.ID(),
			}
			salt, err := user.PasswordSalt()
			if err != nil {
				log.Debug().Err(err).Msg("failed to get salt")
				return ErrorInternalServerError
			}
			hash, err := user.PasswordHash()
			if err != nil {
				log.Debug().Err(err).Msg("failed to get hash")
				return ErrorInternalServerError
			}
			ok, err := srv.verifyPassword(password, salt, hash)
			if err != nil {
				log.Error().Err(err).Msg("failed to verify password")
				return ErrorInternalServerError
			}
			if !ok {
				log.Error().Msg("password is not valid")
				return ErrorUsernameOrPasswordNotValid
			}
			if err := srv.createAuthCookie(w, auth); err != nil {
				log.Error().Err(err).Msg("failed to create 'auth' cookie")
				return ErrorInternalServerError
			}
			return ""
		}()
		if Data.ErrorMessage == "" {
			// the user is redirected to the main page if the authorization data is correct
			log.Debug().Str("redirect", data.EndpointIndex).Msg("success POST HandleLogin")
			redirect(data.EndpointIndex)
			return
		}
	}

	// render of page
	args.Data = Data
	srv.executeTemplate(w, "page_login", args)
}

// HandleLogout is a handler of logout.
func (srv *Service) HandleLogout(w http.ResponseWriter, r *http.Request) {
	a := getAuth(r)
	deleteAuthCookie(w)
	log.Debug().Str("redirect", data.EndpointIndex).Int64("id", a.ID).Msg("client logged out")
	http.Redirect(w, r, data.EndpointIndex, http.StatusSeeOther)
}

func (srv *Service) HandleSignup(w http.ResponseWriter, r *http.Request) {
	redirect := func(endpoint string) {
		http.Redirect(w, r, endpoint, http.StatusSeeOther)
	}
	auth := getAuth(r)

	if auth.IsAuthorized {
		log.Debug().Str("redirect", data.EndpointIndex).Int64("id", auth.ID).Msg("client already signed up")
		redirect(data.EndpointIndex)
		return
	}

	var Data struct {
		ErrorMessage string
	}

	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("signup"),
		},
	}

	// POST method processes data from the user
	if r.Method == http.MethodPost {
		Data.ErrorMessage = func() string {
			if err := r.ParseForm(); err != nil {
				log.Warn().Err(err).Msg("failed to parse form")
				return ErrorBadRequest
			}
			username, password, repeatPassword := r.Form.Get("_username"), r.Form.Get("_password"), r.Form.Get("_repeat_password")
			if err := ValidateUsername(username); err != nil {
				log.Debug().Err(err).Send()
				if fErr, ok := err.(ErrorFrontend); ok {
					return fErr.Frontend
				}
				return ErrorUsernameIsNotValid
			}
			if err := ValidatePassword(password); err != nil {
				log.Debug().Err(err).Send()
				if fErr, ok := err.(ErrorFrontend); ok {
					return fErr.Frontend
				}
				return ErrorPasswordIsNotValid
			}
			if password != repeatPassword {
				log.Debug().Msg("passwords must match")
				return ErrorPasswordsMustMatch
			}
			if strings.ToLower(username) == strings.ToLower(password) {
				log.Debug().Msg("username and password same")
				return ErrorUsernamePasswordSame
			}
			if isVacant, err := model.IsUsernameVacant(srv.redis, username); err != nil {
				log.Error().Err(err).Msg("failed to check if username is vacant")
				return ErrorInternalServerError
			} else if !isVacant {
				log.Debug().Msg("username is not vacant")
				return ErrorUsernameAlreadyTaken
			}
			salt := generatePasswordSalt()
			hash, err := srv.hashPassword(password, salt)
			if err != nil {
				log.Error().Err(err).Msg("failed to hash password")
				return ErrorInternalServerError
			}
			user, err := model.NewUser(srv.redis, username)
			if err != nil {
				log.Error().Err(err).Msg("failed to create user")
				return ErrorInternalServerError
			}
			auth = &data.Auth{
				IsAuthorized: true,
				ID:           user.ID(),
			}
			if err := user.SetPasswordSalt(salt); err != nil {
				log.Error().Err(err).Msg("failed to set salt")
				return ErrorInternalServerError
			}
			if err := user.SetPasswordHash(hash); err != nil {
				log.Error().Err(err).Msg("failed to set hash")
				return ErrorInternalServerError
			}
			if err := srv.createAuthCookie(w, auth); err != nil {
				log.Error().Err(err).Msg("failed to create 'auth' cookie")
				return ErrorInternalServerError
			}
			return ""
		}()
		if Data.ErrorMessage == "" {
			// the user is redirected to the main page if the authorization data is correct
			log.Debug().Str("redirect", data.EndpointIndex).Msg("success POST HandleSignup")
			redirect(data.EndpointIndex)
			return
		}
	}

	// render of page
	args.Data = Data
	srv.executeTemplate(w, "page_signup", args)
}
