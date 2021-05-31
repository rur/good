package example

import (
	"github.com/rur/good/_baseline/site/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	example := treetop.NewView(
		"page/templates/base.html.tmpl",
		hlp.BindEnv(page.BaseHandler),
	)
	placeholder := example.NewDefaultSubView(
		"content",
		"page/example/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)

	hlp.HandleGET("/example",
		exec.NewViewHandler(placeholder))

}
