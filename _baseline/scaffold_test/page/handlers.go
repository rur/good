package page

import (
	"net/http"

	"github.com/rur/good/baseline/scaffold_test/site"
	"github.com/rur/treetop"
)

// ExampleSharedHandler can be used as a handler by multiple pages
func ExampleSharedHandler(env *site.Env, rsp treetop.Response, req *http.Request) interface{} {
	return "Example"
}
