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
			Handler:  "mypageHandler",
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
							Blocks: []routemap.TemplateBlock{
								{
									Name: "form",
									Views: []routemap.RouteView{
										{
											Ref:      "my-form",
											Template: "page/mypage/templates/content/form/my-form.html.tmpl",
											Handler:  "myFormHandler",
											Doc:      "A content form",
											Path:     "/my-page/form",
										},
									},
								},
							},
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

	gotEntries, gotRoutes, err := TemplateDataFromRoutes(def)

	wantEntries := []Entry{
		{
			Assignment: "mypage",
			Block:      "",
			Extends:    "",
			Template:   "page/mypage/templates/mypage.html.tmpl",
			Handler:    "mypageHandler",
			Type:       "PageView",
		},
		{
			Type:    "Spacer",
			Comment: "[[content]]",
		},
		{
			Assignment: "myContent",
			Extends:    "mypage",
			Block:      "content",
			Template:   "page/mypage/templates/content/my-content.html.tmpl",
			Handler:    "myContentHandler",
			Type:       "DefaultSubView",
		},
		{
			Type:    "Spacer",
			Comment: "[[content.form]]",
		},
		{
			Assignment: "myForm",
			Extends:    "myContent",
			Block:      "form",
			Template:   "page/mypage/templates/content/form/my-form.html.tmpl",
			Handler:    "myFormHandler",
			Type:       "SubView",
		},
		{
			Type:    "Spacer",
			Comment: "[[content]]",
		},
		{
			Assignment: "otherContent",
			Extends:    "mypage",
			Block:      "content",
			Template:   "page/mypage/templates/content/other-content.html.tmpl",
			Handler:    "otherContentHandler",
			Type:       "SubView",
		},
		{
			Type:    "Spacer",
			Comment: "[[nav]]",
		},
		{
			Assignment: "",
			Extends:    "mypage",
			Block:      "nav",
			Template:   "page/mypage/templates/nav/my-nav.html.tmpl",
			Handler:    "myNavHandler",
			Type:       "DefaultSubView",
		},
	}
	wantRoutes := []Route{}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v,\n\n want %v", gotEntries, wantEntries)
	}
	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v,\n\n want %v", gotRoutes, wantRoutes)
	}
}
