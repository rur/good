package routemap

import (
	"reflect"
	"testing"

	"github.com/rur/good/generate"
)

func TestTemplateDataFromRoutes(t *testing.T) {
	def := PageRoutes{
		URI: "/my-page",
		RouteView: RouteView{
			Ref:      "mypage",
			Template: "page/mypage/templates/mypage.html.tmpl",
			Handler:  "mypageHandler",
			Doc:      "Test page docs",
			Blocks: []TemplateBlock{
				{
					Name: "content",
					Views: []RouteView{
						{
							Ref:      "my-content",
							Template: "page/mypage/templates/content/my-content.html.tmpl",
							Handler:  "myContentHandler",
							Doc:      "The default content",
							Default:  true,
							Path:     "/my-page",
							Partial:  true,
							Method:   "GET",
							Includes: []string{"my-nav"},
							Blocks: []TemplateBlock{
								{
									Name: "form",
									Views: []RouteView{
										{
											Ref:      "my-form",
											Template: "page/mypage/templates/content/form/my-form.html.tmpl",
											Handler:  "myFormHandler",
											Doc:      "A content form",
											Path:     "/my-page/form",
											Fragment: true,
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
					Views: []RouteView{
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

	// TODO: test missing templates and handlers
	gotEntries, gotRoutes, _, _, err := TemplateDataForRoutes(def, nil, nil)

	wantEntries := []generate.Entry{
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
	wantRoutes := []generate.Route{
		{
			Method:    "GET",
			Path:      "/my-page",
			Includes:  []string{"myNav"},
			Reference: "myContent",
		},
		{
			Path:         "/my-page/form",
			Reference:    "myForm",
			FragmentOnly: true,
		},
		{
			Method:    "POST",
			Path:      "/my-page/other",
			Reference: "otherContent",
			PageOnly:  true,
		},
	}

	if err != nil {
		t.Errorf("TemlateDataFromRoutes() error = %v", err)
		return
	}

	if !reflect.DeepEqual(gotEntries, wantEntries) {
		t.Errorf("TemlateDataFromRoutes() gotEntries = %v,\n\n want %v\n\n", gotEntries, wantEntries)
	}
	if !reflect.DeepEqual(gotRoutes, wantRoutes) {
		t.Errorf("TemlateDataFromRoutes() gotRoutes = %v,\n\n want %v\n\n", gotRoutes, wantRoutes)
	}
}
