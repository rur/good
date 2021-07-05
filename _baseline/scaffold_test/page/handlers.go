package page

import (
	"net/http"

	"github.com/rur/good/baseline/scaffold_test/service"
	"github.com/rur/treetop"
)

// GetBaseHandler Loads handlers data for nav, content and scripts blocks
func GetBaseHandler(title string) ViewHandlerWithEnv {
	return func(env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
		return struct {
			PageTitle string
			Nav       interface{}
			Content   interface{}
			Scripts   interface{}
		}{
			PageTitle: title,
			Nav:       rsp.HandleSubView("nav", req),
			Content:   rsp.HandleSubView("content", req),
			Scripts:   rsp.HandleSubView("scripts", req),
		}
	}
}
