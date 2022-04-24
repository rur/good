package newpage

import (
	"net/http"

	"github.com/rur/good/baseline/starter_test/site"
	"github.com/rur/treetop"
)

// exampleHandler demonstrates a handler with page resources and the site environment
func exampleHandler(rsc *resources, env *site.Env, rsp treetop.Response, req *http.Request) interface{} {
	return struct {
		Description string
	}{
		Description: "Example of a handler, generated for the newpage page",
	}
}
