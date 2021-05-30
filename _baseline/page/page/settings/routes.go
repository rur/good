package settings

import (
	"github.com/rur/treetop"
	"github.com/rur/good/_baseline/site/page"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {
	
	settings := treetop.NewView(
		"page/templates/base.html.tmpl",
		hlp.BindEnv(page.BaseHandler),
	)
	placeholder := settings.NewDefaultSubView(
		"content",
		"page/settings/templates/content/placeholder.html.tmpl",
		hlp.BindEnv(bindResources(placeholderHandler)),
	)
	
	hlp.HandleGET("/settings",
		exec.NewViewHandler(placeholder))
	
}
