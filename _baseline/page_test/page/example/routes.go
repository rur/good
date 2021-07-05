package example

import (
	"github.com/rur/good/baseline/page_test/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Placeholder setup, run `go generate github.com/rur/good/baseline/page_test/page/example` see the starter page

	readme := treetop.NewView(
		"page/example/templates/PLACEHOLDER.html.tmpl",
		treetop.Noop,
	)

	hlp.Handle("/example",
		exec.NewViewHandler(readme).PageOnly())
}
