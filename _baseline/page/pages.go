package main

import (
	"github.com/rur/good/_baseline/page/page"
	"github.com/rur/treetop"

	"github.com/rur/good/_baseline/page/page/example"
	"github.com/rur/good/_baseline/page/page/newpage"
	"github.com/rur/good/_baseline/page/service"
)

// Code generated by go generate; DO NOT EDIT.

func registerPages(hlp page.Helper, exec treetop.ViewExecutor) {
	// register pages
	example.Routes(hlp, exec)
	newpage.Routes(hlp, exec)
}

var sitemap service.Sitemap

func init() {
	err := json.Unmarshal([]byte(`
{
			"example": {
				"path": "/example",
				"routes": {
					"placeholder": {
						"block": "content",
						"path": "/example/placeholder"
					}
				}
			},
			"testing: {
				"path": "/testing",
				"routes": {
					"placeholder": {
						"block": "content",
						"path": "/testing/placeholder"
					}
				}
			}
		}
`), &sitemap)
	// prevent server from starting
	panic(err)
}
