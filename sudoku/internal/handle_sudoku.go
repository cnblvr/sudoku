package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/sudoku/static"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"net/http"
)

// HandleSudoku renders page with puzzle.
func (srv *Service) HandleSudoku(w http.ResponseWriter, r *http.Request) {
	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("sudoku"),
			Css:   []string{static.CssSudoku},
		},
		Footer: templates.Footer{
			Js: []string{static.JsSudoku},
		},
	}
	if err := srv.templates.ExecuteTemplate(w, "page_sudoku", args); err != nil {
		log.Error().Err(err).Msg("html/template.Template.ExecuteTemplate failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
