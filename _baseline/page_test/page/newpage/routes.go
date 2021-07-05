package newpage

import (
	"github.com/rur/good/baseline/page_test/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Placeholder setup, run `go generate github.com/rur/good/baseline/page_test/page/newpage` see the starter page

	readme := treetop.NewView(
		"page/newpage/templates/PLACEHOLDER.html.tmpl",
		treetop.Noop,
	)

	hlp.Handle("/newpage",
		exec.NewViewHandler(readme).PageOnly())
}
