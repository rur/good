package example

import (
	"github.com/rur/good/baseline/scaffold_test/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Code generated by go generate; DO NOT EDIT.

	example := treetop.NewView(
		"page/example/templates/example.html.tmpl",
		hlp.BindEnv(bindResources(exampleHandler)),
	)

	// [[content]]
	examplePlaceholder := example.NewDefaultSubView(
		"content",
		"page/example/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)

	// [[site-nav]]
	example.NewDefaultSubView(
		"site-nav",
		"page/templates/nav.html.tmpl",
		hlp.BindEnv(page.SiteNavHandler),
	)

	// [[scripts]]
	example.NewDefaultSubView(
		"scripts",
		"page/templates/scripts.html.tmpl",
		treetop.Noop,
	)

	hlp.HandleGET("/example",
		exec.NewViewHandler(examplePlaceholder).PageOnly())

}
