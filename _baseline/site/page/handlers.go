package page

import (
	"net/http"

	"github.com/rur/good/_baseline/site/service"
	"github.com/rur/treetop"
)

type link struct {
	Title string
	Path  string
}

// SiteNavHandler is the root handler use for most pages
func SiteNavHandler(env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	var data struct {
		Links []link
	}
	return data
}
