package example

import (
	"github.com/rur/good/_baseline/site/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	example := treetop.NewView(
		"page/example/templates/example.html.tmpl",
		hlp.BindEnv(bindResources(exampleHandler)),
	)
	example.NewDefaultSubView(
		"site-nav",
		"page/templates/nav.html.tmpl",
		hlp.BindEnv(page.SiteNavHandler),
	)
	placeholder := example.NewDefaultSubView(
		"content",
		"page/example/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)
	example.NewDefaultSubView(
		"scripts",
		"page/templates/scripts.html.tmpl",
		treetop.Noop,
	)

	hlp.HandleGET("/example",
		exec.NewViewHandler(placeholder))

}
