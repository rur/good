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
			Template: filepath.Join(templatePath, name+".html.tmpl"),
			Handler:  fmt.Sprintf("hlp.BindEnv(bindResources(%sHandler))", name),
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
							Template: filepath.Join(templatePath, "content", "placeholder.html.tmpl"),
							Handler:  "hlp.BindEnv(bindResources(placeholderHandler))",
						},
					},
				},
				{
					Name: "site-nav",
					Views: []RouteView{
						{
							Ref:      "site-nav",
							Default:  true,
							Template: "page/templates/nav.html.tmpl",
							Handler:  "hlp.BindEnv(page.SiteNavHandler)",
						},
					},
				},
				{
					Name: "scripts",
					Views: []RouteView{
						{
							Ref:      "site-script",
							Default:  true,
							Template: "page/templates/scripts.html.tmpl",
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
