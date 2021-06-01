package newpage

import (
	"net/http"

	"github.com/rur/good/_baseline/page/service"
	"github.com/rur/treetop"
)

// -------------------------
// newpage Handlers
// -------------------------

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
