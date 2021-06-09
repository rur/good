package trivial

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// trivial Handlers
// -------------------------

// Ref: trivial-page
// Doc: Just a single HTML page
func trivialPageHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "trivial Page trivialPageHandler",
	}
	return data
}
