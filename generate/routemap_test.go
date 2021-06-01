package generate

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestLoadRouteRoutemap(t *testing.T) {
	toml, err := ioutil.ReadFile("./testdata/routemap.toml")
	if err != nil {
		t.Fatal("failed ot load test data", err)
	}
	want := &PageRoutes{
		Namespace: "github.com/rur/good/admin/site",
		URI:       "/example",
		RouteView: RouteView{
			Name:     "example",
			Doc:      "Base HTML template for example page",
			Template: "page/templates/base.html.tmpl",
			Handler:  "hlp.BindEnv(page.BaseHandler)",
			Blocks: []TemplateBlock{
				{
					Name: "scripts",
					Views: []RouteView{
						{
							Name:     "site-wide-script",
							Default:  true,
							Template: "page/templates/scripts.html.tmpl",
							Handler:  "treetop.Noop",
						},
					},
				},
				{
					Name: "content",
					Views: []RouteView{
						{
							Name:     "example-placeholder",
							Default:  true,
							Doc:      "Placeholder page",
							Path:     "/example",
							Template: "page/example/templates/content/placedholder.html.tmpl",
							Handler:  "env.Bind(bindResources(placedholderHandler))",
							Blocks: []TemplateBlock{
								{
									Name: "form",
									Views: []RouteView{
										{
											Name:     "placeholder-form",
											Default:  true,
											Doc:      "embedded HTML form",
											Template: "page/example/templates/content/form/placeholderForm.html.tmpl",
											Handler:  "env.Bind(bindResources(placedholderFormHandler))",
										}, {
											Name:     "placeholder-form-preview",
											Fragment: true,
											Method:   "POST",
											Path:     "/example/preview",
											Doc:      "Preview data for submit endpoint",
											Template: "page/example/templates/content/form/placeholderFormPreview.html.tmpl",
											Handler:  "env.Bind(bindResources(placedholderFormPreviewHandler))",
										},
									},
								},
							},
						},
						{
							Name:     "example-submit-endpoint",
							Method:   "POST",
							Page:     true,
							Doc:      "Some form post endpoint",
							Path:     "/example/submit",
							Template: "page/example/templates/content/submit.html.tmpl",
							Handler:  "env.Bind(bindResources(submitHandler))",
						},
					},
				},
			},
		},
	}
	got, err := LoadRouteRoutemap(string(toml))
	if err != nil {
		t.Errorf("LoadRouteRoutemap() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("LoadRouteRoutemap() = %v, want %v", got, want)
	}
}
