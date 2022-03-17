package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/sudoku/static"
	"github.com/cnblvr/sudoku/sudoku/templates"
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
	srv.executeTemplate(w, "page_sudoku", args)
}
