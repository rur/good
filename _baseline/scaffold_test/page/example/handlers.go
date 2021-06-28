package example

import (
	"net/http"

	"github.com/rur/good/baseline/scaffold_test/service"
	"github.com/rur/treetop"
)

// TODO: Delete this after you have run the generate command

// Doc: Handler for the PLACEHOLDER page for the default page boostrap
func placeholderPageHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	data := struct {
		HandlerInfo string
	}{
		HandlerInfo: "PLACEHOLDER handler",
	}
	return data
}
