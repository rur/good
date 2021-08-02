package intro

import (
	"github.com/rur/good/baseline/scaffold_test/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {
	readme := treetop.NewView(
		"page/intro/templates/index.html.tmpl",
		treetop.Noop,
	)

	hlp.Handle("/", exec.NewViewHandler(readme).PageOnly())
}