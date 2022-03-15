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
	const tpl = "page_index"
	if err := srv.templates.ExecuteTemplate(w, tpl, args); err != nil {
		log.Error().Err(err).Str("template", tpl).Msg("failed to execute template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
