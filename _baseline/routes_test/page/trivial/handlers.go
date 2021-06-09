package trivial

import (
	"net/http"

	"github.com/rur/good/baseline/routes_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// trivial Handlers
// -------------------------

// Ref: trivial
// Block: content
// Method: GET
// Doc: Root handler for the trivial page
func trivialHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		SiteNav     interface{}
		Content     interface{}
		Scripts     interface{}
	}{
		HandlerInfo: "trivial Page trivialHandler",
		SiteNav:     rsp.HandleSubView("site-nav", req),
		Content:     rsp.HandleSubView("content", req),
		Scripts:     rsp.HandleSubView("scripts", req),
	}
	return data
}

// Ref: placeholder
// Block: content
// Method: GET
// Doc: This is placeholder content, add your endpoints to the routemap.toml and run go generate
func placeholderHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "trivial Page placeholderHandler",
	}
	return data
}
