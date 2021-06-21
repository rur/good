package example

import (
	"net/http"

	"github.com/rur/good/baseline/scaffold_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// Ref: example-dummy
// Block: content
// Method: GET
// Doc: This is an unused handler for the sake of example
func exampleDummyHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page exampleDummyHandler",
	}
	return data
}
