package sudoku

import (
	"context"
	"github.com/cnblvr/sudoku/data"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

// MiddlewareCookies reads the cookies used in the service and puts them in the context.
func (srv *Service) MiddlewareCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := log.Logger.With().Str("path", r.URL.Path).Logger()
		ctx = context.WithValue(ctx, "logger", log)

		if c, err := r.Cookie("auth"); err == nil {
			a := data.Auth{}
			if err := srv.securecookie.Decode("auth", c.Value, &a); err != nil {
				log.Warn().Err(err).Msg("failed to decode cookie 'auth'")
			} else {
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
		ctx := r.Context()
		auth, log := getAuth(ctx), getLogger(ctx)
		redirect := func(endpoint string) {
			deleteAuthCookie(w)
			http.Redirect(w, r, endpoint, http.StatusSeeOther)
		}

		if !auth.IsAuthorized {
			log.Debug().Str("redirect", data.EndpointIndex).Msg("client is not authorized")
			redirect(data.EndpointIndex)
			return
		}

		_, err := srv.userRepository.GetUserByID(ctx, auth.ID)
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

		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}

func getLogger(ctx context.Context) zerolog.Logger {
	val := ctx.Value("logger")
	logger, ok := val.(zerolog.Logger)
	if !ok {
		return log.Logger
	}
	return logger
}

func getAuth(ctx context.Context) *data.Auth {
	val := ctx.Value("auth")
	if val == nil {
		return &data.Auth{}
	}
	auth, ok := val.(*data.Auth)
	if !ok {
		return &data.Auth{}
	}
	return auth
}
