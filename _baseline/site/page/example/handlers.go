package example

import (
	"net/http"

	"github.com/rur/good/_baseline/site/service"
	"github.com/rur/treetop"
)

// -------------------------
// example Handlers
// -------------------------

// Ref: example// Extends: content
// Method: GET
// Doc: Root handler for the example page
func exampleHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
		SiteNav     interface{}
		Contents    interface{}
		Scripts     interface{}
	}{
		HandlerInfo: "example",
		SiteNav:     rsp.HandleSubView("site-nav", req),
		Contents:    rsp.HandleSubView("contents", req),
		Scripts:     rsp.HandleSubView("scripts", req),
	}
	return data
}

// Ref: placeholder// Extends: content
// Method: GET
// Doc: This is placeholder content, add your endpoints to the routemap.toml and run go generate
func placeholderHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "placeholder",
	}
	return data
}
