package example

import (
	"net/http"

	"github.com/rur/treetop"
	"github.com/rur/good/_baseline/site/app"
)

// -------------------------
// example Handlers
// -------------------------


// placeholder handler DefaultSubView
// Extends: content
// Method: GET
// Doc: This is a placeholder, run go generate command
func placeholderHandler(rsc *resources, env *app.Env, rsp treetop.Response, req *http.Request) interface{} {
	data :=  struct {
		HandlerInfo string
	}{
		HandlerInfo: "placeholder handler",
	}
	return data
}

