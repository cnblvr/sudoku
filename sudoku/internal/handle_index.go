package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HandleIndex is a handler of main page.
func (srv *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	auth := getAuth(r)
	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("index"),
		},
		Auth: auth,
	}
	if auth.IsAuthorized {
		var err error
		args.User, err = srv.userRepository.GetUserByID(ctx, auth.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user")
		}
	}

	srv.executeTemplate(w, "page_index", args)
}
