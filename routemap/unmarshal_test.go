package routemap

import (
	"reflect"
	"testing"

	"github.com/pelletier/go-toml"
)

func TestProcessRoutemapBasic(t *testing.T) {
	tree, err := toml.LoadFile("./testdata/routemap.toml")
	if err != nil {
		t.Fatal("failed to load TOML tree", err)
	}

	want := PageRoutes{
		URI: "/example",
		RouteView: RouteView{
			Ref:      "example",
			Doc:      "Base HTML template for example page",
			Template: "page/templates/base.html.tmpl",
			Handler:  "hlp.BindEnv(page.BaseHandler)",
			Blocks: []TemplateBlock{
				{
					Name: "content",
					Views: []RouteView{
						{
							Ref:      "example-placeholder",
							Default:  true,
							Partial:  true,
							Doc:      "Placeholder page",
							Path:     "/example",
							Template: "page/example/templates/content/placeholder.html.tmpl",
							Handler:  "env.Bind(bindResources(placeholderHandler))",
							Includes: []string{"page-nav"},
							Blocks: []TemplateBlock{
								{
									Name: "form",
									Views: []RouteView{
										{
											Ref:      "placeholder-form",
											Default:  true,
											Doc:      "embedded HTML form",
											Template: "page/example/templates/content/form/placeholderForm.html.tmpl",
											Handler:  "env.Bind(bindResources(placeholderFormHandler))",
										}, {
											Ref:      "placeholder-form-preview",
											Fragment: true,
											Method:   "POST",
											Path:     "/example/preview",
											Doc:      "Preview data for submit endpoint",
											Template: "page/example/templates/content/form/placeholderFormPreview.html.tmpl",
											Handler:  "env.Bind(bindResources(placeholderFormPreviewHandler))",
										},
									},
								},
							},
						},
						{
							Ref:      "example-submit-endpoint",
							Method:   "POST",
							Doc:      "Some form post endpoint",
							Path:     "/example/submit",
							Template: "page/example/templates/content/submit.html.tmpl",
							Handler:  "env.Bind(bindResources(submitHandler))",
						},
					},
				},
				{
					Name: "nav",
					Views: []RouteView{
						{
							Ref:      "page-nav",
							Default:  true,
							Template: "page/example/templates/nav/page-nav.html.tmpl",
							Handler:  "env.Bind(bindResources(pageNavHandler))",
						},
					},
				},
				{
					Name: "scripts",
					Views: []RouteView{
						{
							Ref:      "site-wide-script",
							Default:  true,
							Template: "page/templates/scripts.html.tmpl",
							Handler:  "treetop.Noop",
						},
					},
				},
			},
		},
	}
	got, missT, missH, err := ProcessRoutemap(tree, "page/example")
	if err != nil {
		t.Errorf("GetPageRoutes() error = %v", err)
		return
	}
	if len(missT) > 0 {
		t.Errorf("Unexpected missing templates %v", missT)
		return
	}
	if len(missH) > 0 {
		t.Errorf("Unexpected missing handlers %v", missH)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetPageRoutes() = %v, want %v", got, want)
	}
}
