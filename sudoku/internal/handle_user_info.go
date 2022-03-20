package sudoku

import (
	"fmt"
	"github.com/cnblvr/sudoku/data"
	"github.com/cnblvr/sudoku/sudoku/templates"
	"github.com/rs/zerolog/log"
	"net/http"
)

func (srv *Service) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	redirect := func(endpoint string) {
		http.Redirect(w, r, endpoint, http.StatusSeeOther)
	}
	user := getUser(r)

	var d struct {
		ErrorMessage string
	}

	info, err := user.UserInfo()
	if err != nil {
		log.Error().Err(err).Msg("failed to get user info")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// POST method processes data from the user
	if r.Method == http.MethodPost {
		d.ErrorMessage = func() string {
			if err := r.ParseForm(); err != nil {
				log.Warn().Err(err).Msg("failed to parse form")
				return ErrorBadRequest
			}
			info.Name = r.Form.Get("_name")
			if err := user.SetUserInfo(info); err != nil {
				log.Error().Err(err).Msg("failed to set user info")
				return ErrorInternalServerError
			}
			return ""
		}()
		if d.ErrorMessage == "" {
			log.Debug().Str("redirect", data.EndpointUserInfo).Msg("success POST HandleUserInfo")
			redirect(data.EndpointUserInfo)
			return
		}
	}

	// render of page
	args := templates.Args{
		Header: templates.Header{
			Title: fmt.Sprintf("user's info"),
		},
		User: info,
		Data: d,
	}
	srv.executeTemplate(w, "page_user_info", args)
}
