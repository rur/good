package routemap

import (
	"fmt"
	"path/filepath"

	"github.com/rur/good/generate"
)

// placeholderRoutesConfig will return the default built in routes scaffold for new pages
func PlaceholderRoutesConfig(name, templatePath string) (entries []generate.Entry, routes []generate.Route) {
	var err error
	config := PageRoutes{
		URI: "/example",
		RouteView: RouteView{
			Ref:      name,
			Doc:      fmt.Sprintf("Base HTML template for %s page", name),
			Template: filepath.Join("page", "templates", "base.html.tmpl"),
			Handler:  fmt.Sprintf("hlp.BindEnv(page.GetBaseHandler(%#v))", name+" Page"),
			Blocks: []TemplateBlock{
				{
					Name: "content",
					Views: []RouteView{
						{
							Ref:      name + "-placeholder",
							Default:  true,
							Method:   "GET",
							Doc:      "Placeholder page",
							Path:     "/" + name,
							Template: filepath.Join(templatePath, "placeholder.html.tmpl"),
							Handler:  fmt.Sprintf("treetop.Constant(%#v)", name),
						},
					},
				},
				{
					Name: "nav",
					Views: []RouteView{
						{
							Ref:      "page-nav",
							Default:  true,
							Template: "::empty::",
							Handler:  "treetop.Noop",
						},
					},
				},
				{
					Name: "scripts",
					Views: []RouteView{
						{
							Ref:      "site-script",
							Default:  true,
							Template: "::empty::",
							Handler:  "treetop.Noop",
						},
					},
				},
			},
		},
	}
	// TODO: add templates and routes for placeholder
	entries, routes, _, _, err = TemplateDataForRoutes(config, nil, nil)
	if err != nil {
		panic(err)
	}
	return
}
