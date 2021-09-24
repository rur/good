package trivial

import (
	"github.com/rur/good/baseline/routes_test/page"
	"github.com/rur/treetop"
)

// Routes is the plumbing code for page endpoints, templates and handlers
func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Code created by go generate. You should edit the routemap.toml file; DO NOT EDIT.

	trivialPage := treetop.NewView(
		"page/trivial/templates/trivial-page.html.tmpl",
		hlp.BindEnv(bindResources(trivialPageHandler)),
	)

	hlp.Handle("/trivial",
		exec.NewViewHandler(trivialPage).PageOnly())

}
