package templates

import (
	"context"
	"embed"
	"github.com/cnblvr/sudoku/data"
	"html/template"
	"sort"
)

// Templates is a repository of HTML template files.
//go:embed common/*.gohtml *.gohtml
var Templates embed.FS

// Common returns an array of common templates.
func Common() []string {
	return []string{"common/header.gohtml", "common/footer.gohtml"}
}

var CommonFunctions = template.FuncMap{
	"sort_navigation": func(nav []Navigation) []Navigation {
		sort.Slice(nav, func(i, j int) bool {
			return nav[i].Weight < nav[j].Weight
		})
		return nav
	},
	"add_internal_css": func(h Header, css ...template.CSS) Header {
		h.CssInternal = append(h.CssInternal, css...)
		return h
	},
}

// Args are arguments for executing templates.
type Args struct {
	Header Header
	Auth   *data.Auth
	Data   interface{}
	Footer Footer
}

func NewHeader(ctx context.Context, h Header) Header {
	if h.Title == "" {
		h.Title = "unknown page"
	}
	h.Navigation = append(h.Navigation,
		Navigation{Label: "Home", Path: data.EndpointIndex, Weight: 0},
	)

	auth := data.GetCtxAuth(ctx)
	if auth.IsAuthorized {
		h.Navigation = append(h.Navigation, Navigation{
			Label:  "Log out",
			Path:   data.EndpointLogout,
			Weight: 999,
		})
	} else {
		h.Navigation = append(h.Navigation, Navigation{
			Label:  "Log in",
			Path:   data.EndpointLogin,
			Weight: 979,
		})
		h.Navigation = append(h.Navigation, Navigation{
			Label:  "Sign up",
			Path:   data.EndpointSignup,
			Weight: 989,
		})
	}

	return h
}

// Header is placed in the common 'header' template.
type Header struct {
	// Title of page.
	Title      string
	Navigation []Navigation
	// List of css that are used on the page.
	Css         []string
	CssInternal []template.CSS
	// List of js that are used on the page.
	Js []string // todo move to footer
}

type Navigation struct {
	Label  string
	Path   string
	Weight int
}

// Footer is placed in the common 'footer' template.
type Footer struct {
}
