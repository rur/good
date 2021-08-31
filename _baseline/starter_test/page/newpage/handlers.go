package newpage

import (
	"net/http"

	"github.com/rur/good/baseline/starter_test/service"
	"github.com/rur/treetop"
)

// baseHandler is the default top level handler for the newpage page
func baseHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		PageTitle string
		Nav       interface{}
		Content   interface{}
		Styles    interface{}
		Scripts   interface{}
	}{
		PageTitle: "newpage Page",
		Nav:       rsp.HandleSubView("nav", req),
		Content:   rsp.HandleSubView("content", req),
		Styles:    rsp.HandleSubView("styles", req),
		Scripts:   rsp.HandleSubView("scripts", req),
	}
	return data
}
