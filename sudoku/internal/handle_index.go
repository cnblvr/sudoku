package sudoku

import (
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"html/template"
	"net/http"
)

// HandleIndex renders index page.
func (srv *Service) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// Temporary parse templates
	t, err := template.ParseFS(templates.Templates, append(templates.Common(), "index.gohtml")...)
	if err != nil {
		log.Error().Err(err).Msg("html/template.ParseFS failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := t.ExecuteTemplate(w, "index", struct{}{}); err != nil {
		log.Error().Err(err).Msg("html/template.Template.ExecuteTemplate failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
