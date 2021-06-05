package routemap

import (
	"reflect"
	"strings"
	"testing"

	toml "github.com/pelletier/go-toml"
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

func TestProcessRoutemapErrors(t *testing.T) {
	type args struct {
		toml         string
		templatePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr string
	}{
		{
			name: "bad ref",
			args: args{
				toml: `
_ref = "invalid ref"
				`,
			},
			wantErr: ":1:1: Unknown or invalid _ref 'invalid ref', references be all lowercase joined by a dash '-'",
		},
		{
			name: "bad ref 2",
			args: args{
				toml: `
_ref = "invalidREF"
				`,
			},
			wantErr: ":1:1: Unknown or invalid _ref 'invalidREF', references be all lowercase joined by a dash '-'",
		},
		{
			name: "bad block name",
			args: args{
				toml: `
_ref = "okref"
	[[testing_fail]]
	 _ref = "testing-block"
				`,
			},
			wantErr: ":3:2: Unknown or invalid key 'testing_fail', block names must be all lowercase joined by a dash '-'",
		},
		{
			name: "duplicate ref",
			args: args{
				toml: `
_ref = "my-ref"
	[[testing]]
	 _ref = "my-ref"
				`,
			},
			wantErr: ":3:2: duplicate _ref 'my-ref', already used in routemap file at line 2, column 1",
		},
		{
			name: "duplicate ref",
			args: args{
				toml: `
_ref = "mypage"
	[[testing]]
	 _ref = "my-sub-view"

	 unknown = "something wrong here"
				`,
			},
			wantErr: `:6:3: invalid value for key 'unknown', expecting an array of tables, got "something wrong here"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := toml.Load(tt.args.toml)
			if err != nil {
				t.Fatal("Bad toml string", tt.name, err.Error())
			}
			_, _, _, err = ProcessRoutemap(tree, tt.args.templatePath)
			if err == nil {
				t.Error("ProcessRoutemap() error handling, expecting an error, got none")
				return
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("ProcessRoutemap() expecting error to contain %v, got %v", tt.wantErr, err.Error())
			}
		})
	}
}
