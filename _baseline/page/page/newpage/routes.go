package newpage

import (
	"github.com/rur/good/_baseline/page/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Code generated by go generate; DO NOT EDIT.

	newpage := treetop.NewView(
		"page/newpage/templates/newpage.html.tmpl",
		hlp.BindEnv(bindResources(newpageHandler)),
	)

	// [content]
	placeholder := newpage.NewDefaultSubView(
		"content",
		"page/newpage/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)

	// [scripts]
	newpage.NewDefaultSubView(
		"scripts",
		"page/templates/scripts.html.tmpl",
		treetop.Noop,
	)

	// [site-nav]
	newpage.NewDefaultSubView(
		"site-nav",
		"page/templates/nav.html.tmpl",
		hlp.BindEnv(page.SiteNavHandler),
	)

	hlp.HandleGET("/newpage",
		exec.NewViewHandler(placeholder))

}
