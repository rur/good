package generate

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetRoutesFile(t *testing.T) {
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
						},
						{
							Name:     "example-submit-endpoint",
							Method:   "POST",
							Fragment: true,
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
	got, err := GetRoutes(string(toml))
	if err != nil {
		t.Errorf("GetRoutes() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetRoutes() = %v, want %v", got, want)
	}
}
