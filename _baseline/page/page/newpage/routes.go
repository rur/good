package newpage

import (
	"github.com/rur/good/_baseline/page/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	newpage := treetop.NewView(
		"page/newpage/templates/newpage.html.tmpl",
		hlp.BindEnv(bindResources(newpageHandler)),
	)
	newpage.NewDefaultSubView(
		"site-nav",
		"page/templates/nav.html.tmpl",
		hlp.BindEnv(page.SiteNavHandler),
	)
	placeholder := newpage.NewDefaultSubView(
		"content",
		"page/newpage/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)
	newpage.NewDefaultSubView(
		"scripts",
		"page/templates/scripts.html.tmpl",
		treetop.Noop,
	)

	hlp.HandleGET("/newpage",
		exec.NewViewHandler(placeholder))

}
