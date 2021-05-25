package generate

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type Handler struct {
	Info       string
	Type       string // "Fragment" "Partial"
	Extends    string
	Method     string
	Doc        string
	Identifier string
	Blocks     []Block
}

type Block struct {
	FieldName string
	Name      string
}

type Entry struct {
	Type       string // "SubView" "DefaultSubView" "Spacer"
	Assignment string
	Extends    string
	Block      string
	Template   string
	Handler    string
	Name       string
}

type Route struct {
	Method    string
	Path      string
	Includes  []string
	Reference string
	Type      string // "Page" "Fragment" ""
}

// ScaffoldPage will assemble files for adding a new page to the site scaffold
func ScaffoldPage(siteModule, siteDir, name string, scaffold fs.FS) (files []File, err error) {
	data := struct {
		Name      string // Go package name for page
		Namespace string
		Handlers  []Handler
		PageEntry Entry
		Entries   []Entry
		Routes    []Route
		Templates string
		PagePath  string
	}{
		PagePath:  strings.Join([]string{siteModule, "page", name}, "/"),
		Name:      name,
		Namespace: siteModule,
		Handlers: []Handler{
			{
				Info:       "placeholder handler",
				Type:       "DefaultSubView",
				Extends:    "content",
				Method:     "GET",
				Doc:        "This is a placeholder, run go generate command",
				Identifier: "placeholderHandler",
			},
		},
		PageEntry: Entry{
			Assignment: name,
			Template:   filepath.Join(siteDir, "page", "templates", "base.html.tmpl"),
			Handler:    "hlp.BindEnv(page.BaseHandler)",
		},
		Entries: []Entry{{
			Assignment: "placeholder",
			Type:       "DefaultSubView",
			Extends:    name,
			Template:   filepath.Join(siteDir, "page", name, "templates", "placeholder.html.tmpl"),
			Handler:    "hlp.BindEnv(bindResources(placeholderHandler))",
		}},
		Routes: []Route{{
			Method:    "GET",
			Path:      "/" + name,
			Reference: "placeholder",
		}},
		Templates: filepath.Join(siteDir, "page", name, "templates"),
	}

	pageDir := filepath.Join(siteDir, "page", name)

	// page/name/gen.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "gen.go",
		Contents: mustExecute("scaffold/page/name/gen.go.tmpl", data, scaffold),
	})
	// page/name/handlers.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "handles.go",
		Contents: mustExecute("scaffold/page/name/handlers.go.tmpl", data, scaffold),
	})
	// page/name/resources.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "resources.go",
		Contents: mustExecute("scaffold/page/name/resources.go.tmpl", data, scaffold),
	})
	// page/name/routemap.toml
	files = append(files, File{
		Dir:      pageDir,
		Name:     "routemap.toml",
		Contents: mustExecute("scaffold/page/name/routemap.toml.tmpl", data, scaffold),
	})
	// page/name/routes.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "routes.go",
		Contents: mustExecute("scaffold/page/name/routes.go.tmpl", data, scaffold),
	})
	// page/name/templates/content/placeholder.html.tmpl
	files = append(files, File{
		Dir:      filepath.Join(pageDir, "templates", "content"),
		Name:     "placeholder.html.tmpl",
		Contents: mustExecute("scaffold/page/name/templates/content/placeholder.html.tmpl.tmpl", data, scaffold),
	})
	return
}
