package newpage

import (
	"github.com/rur/good/_baseline/page/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	newpage := treetop.NewView(
		"page/templates/base.html.tmpl",
		hlp.BindEnv(page.BaseHandler),
	)
	placeholder := newpage.NewDefaultSubView(
		"content",
		"page/newpage/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)

	hlp.HandleGET("/newpage",
		exec.NewViewHandler(placeholder))

}
