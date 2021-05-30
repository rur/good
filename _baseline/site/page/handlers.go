package page

import (
	"net/http"

	"github.com/rur/treetop"
	"github.com/rur/good/_baseline/site/service"
)

// BaseHandler is the root handler use for most pages
func BaseHandler(env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	return struct {
		PageTitle string
		Content   interface{}
	}{
		PageTitle: "Base page",
		Content:   rsp.HandleSubView("content", req),
	}
}
