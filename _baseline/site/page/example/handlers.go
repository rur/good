package example

import (
	"net/http"

	"github.com/rur/good/_baseline/site/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// base handle for example page
// Doc: Root handle for the main page
func exampleHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	return struct {
		PageTitle string
		SiteNav   interface{}
		Content   interface{}
		Scripts   interface{}
	}{
		PageTitle: "Some Example Page",
		SiteNav:   rsp.HandleSubView("site-nav", req),
		Content:   rsp.HandleSubView("content", req),
		Scripts:   rsp.HandleSubView("scripts", req),
	}
}

// placeholder handler DefaultSubView
// Extends: content
// Method: GET
// Doc: This is a placeholder, run go generate command
func placeholderHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "placeholder handler",
	}
	return data
}
