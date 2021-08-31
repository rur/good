package page

import (
	"net/http"

	"github.com/rur/good/baseline/starter_test/service"
	"github.com/rur/treetop"
)

// ExampleSharedHandler can be used as a handler by multiple pages
func ExampleSharedHandler(env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	return "Example"
}
