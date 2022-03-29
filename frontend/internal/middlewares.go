package frontend

import (
	"github.com/cnblvr/sudoku/data"
	"github.com/pkg/errors"
	zlog "github.com/rs/zerolog/log"
	"net/http"
)

// MiddlewareCookies reads the cookies used in the service and puts them in the context.
func (srv *Service) MiddlewareCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = data.SetCtxReqID(ctx, data.GenerateReqID())
		log := zlog.Logger.With().
			Str("path", r.URL.Path).
			Str("reqid", data.GetCtxReqID(ctx)).
			Logger()

		if c, err := r.Cookie("auth"); err == nil {
			auth := data.Auth{}
			if err := srv.securecookie.Decode("auth", c.Value, &auth); err != nil {
				log.Warn().Err(err).Msg("failed to decode cookie 'auth'")
			} else {
				token, err := srv.userRepository.GetToken(ctx, auth.Token)
				if err != nil {
					log.Debug().Err(err).Msg("failed to get token")
					token, err = srv.userRepository.CreateToken(ctx, auth.UserID(), 0)
					if err != nil {
						log.Debug().Err(err).Msg("failed to create token")
					}
				}
				auth.ID = token.UserID
				if auth.ID > 0 {
					auth.IsAuthorized = true
				}
				log = log.With().Int64("user_id", auth.ID).Logger()
				ctx = data.SetCtxAuth(ctx, &auth)
				log.Debug().Interface("auth", auth).Send()
			}
		}

		ctx = data.SetCtxLog(ctx, log)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// MiddlewareMustBeAuthorized does not allow further operations if the user is not logged in.
// MiddlewareCookies pre-middleware required.
func (srv *Service) MiddlewareMustBeAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth, log := data.GetCtxAuth(ctx), data.GetCtxLog(ctx)
		redirect := func(endpoint string) {
			deleteAuthCookie(w)
			http.Redirect(w, r, endpoint, http.StatusSeeOther)
		}

		if !auth.IsAuthorized {
			log.Debug().Str("redirect", data.EndpointIndex).Msg("client is not authorized")
			redirect(data.EndpointIndex)
			return
		}

		user, err := srv.userRepository.GetUserByID(ctx, auth.ID)
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
		ctx = data.SetCtxUser(ctx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
		return
	})
}
