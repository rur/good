package page

import (
	"net/http"

	"github.com/rur/good/_baseline/site/service"
	"github.com/rur/treetop"
)

type Link struct {
	Title string
	Path  string
}

// SiteNavHandler is the root handler use for most pages
func SiteNavHandler(env *service.Env, rsp treetop.Response, req *http.Request) interface{} {
	var data struct {
		Links []Link
	}
	for _, page := range env.Sitemap.Pages() {
		data.Links = append(data.Links, Link{
			Title: page,
			Path:  env.Sitemap[page].URI,
		})
	}
	return data
}
