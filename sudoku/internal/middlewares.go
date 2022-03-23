package sudoku

import (
	"context"
	"github.com/cnblvr/sudoku/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

// MiddlewareCookies reads the cookies used in the service and puts them in the context.
func (srv *Service) MiddlewareCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.With().Str("path", r.URL.Path).Logger()
		ctx := r.Context()
		ctx = context.WithValue(ctx, "auth", &data.Auth{})
		for _, name := range []string{"auth"} {
			c, err := r.Cookie(name)
			if err != nil {
				log.Debug().Err(err).Msg("cookie 'auth' not found")
				continue
			}
			switch c.Name {

			// read 'auth' cookie
			case "auth":
				a := data.Auth{}
				if err := srv.securecookie.Decode("auth", c.Value, &a); err != nil {
					log.Warn().Err(err).Msg("failed to decode cookie 'auth'")
					continue
				}
				ctx = context.WithValue(ctx, "auth", &a)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MiddlewareMustBeAuthorized does not allow further operations if the user is not logged in.
// MiddlewareCookies pre-middleware required.
func (srv *Service) MiddlewareMustBeAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := log.With().Str("path", r.URL.Path).Logger()
		ctx := r.Context()
		redirect := func(endpoint string) {
			deleteAuthCookie(w)
			http.Redirect(w, r, endpoint, http.StatusSeeOther)
		}
		a := getAuth(r)

		if !a.IsAuthorized {
			log.Debug().Str("redirect", data.EndpointIndex).Msg("client is not authorized")
			redirect(data.EndpointIndex)
			return
		}

		user, err := srv.userRepository.GetUserByID(ctx, a.ID)
		if err != nil {
			if errors.Is(err, data.ErrUserNotFound) {
				log.Debug().Str("redirect", data.EndpointLogout).Msg("user not found")
				redirect(data.EndpointLogout)
				return
			}
			log.Debug().Str("redirect", data.EndpointIndex).Msg("failed to get user")
			redirect(data.EndpointIndex)
			return
		}
		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

// getAuth get authorization data from request's context.
func getAuth(r *http.Request) *data.Auth {
	return r.Context().Value("auth").(*data.Auth)
}

func getUser(r *http.Request) *data.User {
	return r.Context().Value("user").(*data.User)
}
