package example

import (
	"github.com/rur/good/baseline/scaffold_test/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Placeholder setup, run `go generate github.com/rur/good/baseline/scaffold_test/page/example` see the starter page

	readme := treetop.NewView(
		"page/example/templates/PLACEHOLDER.html.tmpl",
		hlp.BindEnv(bindResources(placeholderPageHandler)),
	)

	hlp.Handle("/example",
		exec.NewViewHandler(readme).PageOnly())
}
