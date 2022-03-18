package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/model"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HandleIndex is a handler of main page.
func (srv *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	auth := getAuth(r)
	redis := srv.redis.Get()
	defer redis.Close()
	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("index"),
		},
		Auth: auth,
	}
	if auth.IsAuthorized {
		var err error
		user, _, err := model.UserByID(redis, auth.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user")
		}
		args.User, err = user.UserInfo()
		if err != nil {
			log.Error().Err(err).Msg("failed to get user info")
		}
	}

	srv.executeTemplate(w, "page_index", args)
}
