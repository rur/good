package newpage

import (
	"net/http"

	"github.com/rur/good/baseline/page/service"
	"github.com/rur/treetop"
)

// -------------------------
// newpage Handlers
// -------------------------

// Ref: newpage
// Extends: content
// Method: GET
// Doc: Root handler for the newpage page
func newpageHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		SiteNav     interface{}
		Content     interface{}
		Scripts     interface{}
	}{
		HandlerInfo: "newpage Page newpageHandler",
		SiteNav:     rsp.HandleSubView("site-nav", req),
		Content:     rsp.HandleSubView("content", req),
		Scripts:     rsp.HandleSubView("scripts", req),
	}
	return data
}

// Ref: placeholder
// Extends: content
// Method: GET
// Doc: This is placeholder content, add your endpoints to the routemap.toml and run go generate
func placeholderHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "newpage Page placeholderHandler",
	}
	return data
}
