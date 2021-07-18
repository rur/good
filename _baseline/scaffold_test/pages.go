package main

import (
	"github.com/rur/good/baseline/scaffold_test/page"
	"github.com/rur/treetop"

	"github.com/rur/good/baseline/scaffold_test/page/intro"
)

// Code generated by go generate; DO NOT EDIT.

func registerPages(hlp page.Helper, exec treetop.ViewExecutor) {
	// register pages
	intro.Routes(hlp, exec)
}
