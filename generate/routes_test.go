package generate

import (
	"reflect"
	"testing"

	"github.com/rur/good/routemap"
)

func TestTemplateDataFromRoutes(t *testing.T) {
	def := routemap.PageRoutes{
		Namespace: "github.com/rur/example/site/page/mypage",
		URI:       "/my-page",
		RouteView: routemap.RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "hlp.BindEnv(bindResources(mypageHandler))",
			Doc:      "Test page docs",
			Blocks: []routemap.TemplateBlock{
				{
					Name: "content",
					Views: []routemap.RouteView{
						{
							Ref:      "my-content",
							Template: "page/mypage/templates/content/my-content.html.tmpl",
							Handler:  "myContentHandler",
							Doc:      "The default content",
							Default:  true,
							Path:     "/my-page",
							Method:   "GET",
						},
						{
							Ref:      "other-content",
							Template: "page/mypage/templates/content/other-content.html.tmpl",
							Handler:  "otherContentHandler",
							Doc:      "The other content",
							Path:     "/my-page/other",
							Method:   "POST",
						},
					},
				},
				{
					Name: "nav",
					Views: []routemap.RouteView{
						{
							Ref:      "my-nav",
							Template: "page/mypage/templates/nav/my-nav.html.tmpl",
							Handler:  "myNavHandler",
							Doc:      "The default nav",
							Default:  true,
						},
					},
				},
			},
		},
	}

	gotHandlers, gotEntries, gotRoutes, err := TemplateDataFromRoutes(def)

	wantHandlers := []Handler{
		{
			Ref:        "mypage",
			Type:       "PageView",
			Doc:        "Test page docs",
			Identifier: "mypageHandler",
		},
		{
			Ref:        "my-content",
			Extends:    "content",
			Type:       "DefaultSubView",
			Method:     "GET",
			Doc:        "The default content",
			Identifier: "myContentHandler",
		},
		{
			Ref:        "other-content",
			Extends:    "content",
			Type:       "SubView",
			Method:     "POST",
			Doc:        "The other content",
			Identifier: "otherContentHandler",
		},
		{
			Ref:        "my-nav",
			Extends:    "nav",
			Type:       "DefaultSubView",
			Doc:        "The default nav",
			Identifier: "myNavHandler",
		},
	}
	wantEntries := []Entry{}
	wantRoutes := []Route{}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotHandlers, wantHandlers) {
		t.Errorf("TemlateDataFromRoutes() gotHandlers = %v, want %v", gotHandlers, wantHandlers)
	}
	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v, want %v", gotEntries, wantEntries)
	}
	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v, want %v", gotRoutes, wantRoutes)
	}
}
