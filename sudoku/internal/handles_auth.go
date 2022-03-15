package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/sudoku/internal/db"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
)

// HandleLogin is a handler of login page.
func (srv *Service) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	redirect := func(endpoint string) {
		http.Redirect(w, r, endpoint, http.StatusSeeOther)
	}
	auth := getAuth(r)

	// If the user is already logged in, then redirect to the main page.
	if auth.IsAuthorized {
		log.Debug().Str("redirect", data.EndpointIndex).Str("username", auth.Username).Msg("client already logged in")
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
		err := func() error {
			err := r.ParseForm()
			if err != nil {
				Data.ErrorMessage = "Bad Request"
				log.Warn().Err(err).Msg("failed to parse form")
				return err
			}
			username, password := r.Form.Get("_username"), r.Form.Get("_password")
			if err := data.ValidateUsername(username); err != nil {
				Data.ErrorMessage = "username or password not valid"
				return err
			}
			user, isExists, err := db.GetUser(ctx, srv.redis, username)
			if err != nil {
				Data.ErrorMessage = "username or password not valid"
				log.Error().Err(err).Msg("failed to GetUser")
				return err
			}
			if !isExists {
				Data.ErrorMessage = "username or password not valid"
				log.Debug().Err(err).Msg("username is not exists")
				return fmt.Errorf("username is not exists")
			}
			auth = &data.Auth{
				IsAuthorized: true,
				Username:     strings.ToLower(username),
			}
			ok, err := srv.verifyPassword(password, user.PasswordSalt, user.PasswordHash)
			if err != nil {
				Data.ErrorMessage = "Internal Server Error"
				log.Error().Err(err).Msg("failed to verify password")
				return err
			}
			if !ok {
				Data.ErrorMessage = "username or password not valid"
				log.Error().Msg("password is not valid")
				return fmt.Errorf("password is not valid")
			}
			if err := srv.createAuthCookie(w, auth); err != nil {
				Data.ErrorMessage = "Internal Server Error"
				log.Error().Err(err).Msg("failed to create 'auth' cookie")
				return err
			}
			return nil
		}()
		if err == nil {
			// the user is redirected to the main page if the authorization data is correct
			log.Debug().Str("redirect", data.EndpointIndex).Msg("success POST HandleLogin")
			redirect(data.EndpointIndex)
			return
		}
	}

	// render of page
	args.Data = Data
	const tpl = "page_login"
	if err := srv.templates.ExecuteTemplate(w, tpl, args); err != nil {
		log.Error().Err(err).Str("template", tpl).Msg("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// HandleLogout is a handler of logout.
func (srv *Service) HandleLogout(w http.ResponseWriter, r *http.Request) {
	a := getAuth(r)
	deleteAuthCookie(w)
	log.Debug().Str("redirect", data.EndpointIndex).Str("username", a.Username).Msg("client logged out")
	http.Redirect(w, r, data.EndpointIndex, http.StatusSeeOther)
}

func (srv *Service) HandleSignup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	redirect := func(endpoint string) {
		http.Redirect(w, r, endpoint, http.StatusSeeOther)
	}
	auth := getAuth(r)

	if auth.IsAuthorized {
		log.Debug().Str("redirect", data.EndpointIndex).Str("username", auth.Username).Msg("client already signed up")
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
		err := func() error {
			err := r.ParseForm()
			if err != nil {
				Data.ErrorMessage = "Bad Request"
				log.Warn().Err(err).Msg("failed to parse form")
				return err
			}
			username, password := r.Form.Get("_username"), r.Form.Get("_password")
			if err := data.ValidateUsername(username); err != nil {
				Data.ErrorMessage = err.Error()
				return err
			}
			if err := data.ValidatePassword(password); err != nil {
				Data.ErrorMessage = err.Error()
				return err
			}
			if exists, err := db.IsExistsUser(ctx, srv.redis, username); err != nil {
				Data.ErrorMessage = "Internal Server Error"
				log.Error().Err(err).Msg("failed to IsExistsUser")
				return err
			} else if exists {
				Data.ErrorMessage = "username already taken"
				return fmt.Errorf("username already is registered")
			}
			user := data.User{
				Username:     username,
				PasswordSalt: generatePasswordSalt(),
			}
			auth = &data.Auth{
				IsAuthorized: true,
				Username:     strings.ToLower(username),
			}
			user.PasswordHash, err = srv.hashPassword(password, user.PasswordSalt)
			if err != nil {
				Data.ErrorMessage = "Internal Server Error"
				log.Error().Err(err).Msg("failed to hash password")
				return err
			}
			if err := db.CreateUser(ctx, srv.redis, user); err != nil {
				Data.ErrorMessage = "Internal Server Error"
				log.Error().Err(err).Msg("failed to create user")
				return err
			}
			if err := srv.createAuthCookie(w, auth); err != nil {
				Data.ErrorMessage = "Internal Server Error"
				log.Error().Err(err).Msg("failed to create 'auth' cookie")
				return err
			}
			return nil
		}()
		if err == nil {
			// the user is redirected to the main page if the authorization data is correct
			log.Debug().Str("redirect", data.EndpointIndex).Msg("success POST HandleSignup")
			redirect(data.EndpointIndex)
			return
		}
	}

	// render of page
	args.Data = Data
	const tpl = "page_signup"
	if err := srv.templates.ExecuteTemplate(w, tpl, args); err != nil {
		log.Error().Err(err).Str("template", tpl).Msg("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
