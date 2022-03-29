package frontend

import (
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/frontend/templates"
	"net/http"
)

// HandleIndex is a handler of main page.
func (srv *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	d := struct {
		Username string
	}{}

	auth, log := data.GetCtxAuth(ctx), data.GetCtxLog(ctx)
	if auth.IsAuthorized {
		var err error
		user, err := srv.userRepository.GetUserByID(ctx, auth.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get user")
		} else {
			d.Username = user.Username
		}
	}

	args := templates.Args{
		Header: templates.NewHeader(ctx, templates.Header{
			Title: fmt.Sprintf("index"),
		}),
		Data: d,
		Auth: auth,
	}
	srv.executeTemplate(w, "page_index", args)
}
