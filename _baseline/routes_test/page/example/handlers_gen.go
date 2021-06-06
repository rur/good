package example

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// Ref: other-content
// Block: content
// Method:
// Doc: other Content Page
func otherContentHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "example Page otherContentHandler",
	}
	return data
}
