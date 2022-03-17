package templates

import (
	"embed"
	"github.com/cnblvr/sudoku/data"
)

// Templates is a repository of HTML template files.
//go:embed common/*.gohtml *.gohtml
var Templates embed.FS

// Common returns an array of common templates.
func Common() []string {
	return []string{"common/header.gohtml", "common/footer.gohtml"}
}

// Args are arguments for executing templates.
type Args struct {
	Header Header
	Auth   *data.Auth
	User   data.UserInfo
	Data   interface{}
	Footer Footer
}

// Header is placed in the common 'header' template.
type Header struct {
	// Title of page.
	Title string
	// List of css that are used on the page.
	Css []string
}

// Footer is placed in the common 'footer' template.
type Footer struct {
	// List of js that are used on the page.
	Js []string
}
