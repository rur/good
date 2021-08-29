package newpage

import (
	"net/http"

	"github.com/rur/good/baseline/page_test/service"
	"github.com/rur/treetop"
)

// -------------------------
// newpage Handlers
// -------------------------

// Ref: landing
// Block: subsection
func landingHandler(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	return struct {
		SubsectionTabs interface{}
	}{
		SubsectionTabs: rsp.HandleSubView("subsection-tabs", req),
	}
}
