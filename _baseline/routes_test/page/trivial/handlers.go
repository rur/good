package trivial

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// trivialHandler is the default top level handler for the trivial page
func trivialHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		PageTitle string
		Nav       interface{}
		Content   interface{}
		Scripts   interface{}
	}{
		PageTitle: "trivial Page",
		Nav:       rsp.HandleSubView("nav", req),
		Content:   rsp.HandleSubView("content", req),
		Scripts:   rsp.HandleSubView("scripts", req),
	}
	return data
}
