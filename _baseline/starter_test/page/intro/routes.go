package intro

import (
	"github.com/rur/good/baseline/starter_test/page"
	"github.com/rur/treetop"
)

func Routes(hlp page.Helper, exec treetop.ViewExecutor) {

	// Code generated by go generate; DO NOT EDIT.

	index := treetop.NewView(
		"page/intro/templates/index.html.tmpl",
		treetop.Noop,
	)

	// [[diagram]]
	index.NewDefaultSubView(
		"diagram",
		"page/intro/templates/diagram.svg",
		treetop.Noop,
	)

	hlp.Handle("/",
		exec.NewViewHandler(index).PageOnly())

}
