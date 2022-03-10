package templates

import "embed"

// Templates is a repository of HTML template files.
//go:embed common/*.gohtml *.gohtml
var Templates embed.FS

// Common returns an array of common templates.
func Common() []string {
	return []string{"common/header.gohtml", "common/footer.gohtml"}
}
