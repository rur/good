package generate

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

// Entry for the routes.go file
type Entry struct {
	Type       string // "SubView" "DefaultSubView" "Spacer"
	Assignment string
	Extends    string
	Block      string
	Template   string
	Handler    string
	Comment    string
}

// Route is a path mapped to a view definition
type Route struct {
	Method       string
	Path         string
	Includes     []string
	Reference    string
	PageOnly     bool
	FragmentOnly bool
}

// HTMLTemplate is data for a page template document
type HTMLTemplate struct {
	Path   string
	Blocks []HandleBlock
}

// Handler is data for a handler function which should be created
type Handler struct {
	Ref        string
	Block      string
	Method     string
	Doc        string
	Identifier string
	Blocks     []HandleBlock
}

// HandleBlock is the details of sub-views which should
// be delegated to in the handler
type HandleBlock struct {
	FieldName string
	Name      string
}

// RoutesScaffold will generate all files for the good routes command
func RoutesScaffold(
	sitePkg GoPackage,
	pageName string,
	entries []Entry,
	routes []Route,
	templates []HTMLTemplate,
	handlers []Handler,
	scaffold fs.FS,
) (files []File, err error) {
	data := struct {
		Name      string
		Namespace string
		Entries   []Entry
		Routes    []Route
	}{
		Name:      pageName,
		Namespace: sitePkg.ImportPath,
		Entries:   entries,
		Routes:    routes,
	}
	// page/name/routes.go
	files = append(files, File{
		Dir:      filepath.Join("page", pageName),
		Name:     "routes.go",
		Contents: mustExecute("scaffold/page/name/routes.go.tmpl", data, scaffold),
	})
	if len(handlers) > 0 {
		// page/name/handlers.go
		files = append(files, File{
			Dir:      filepath.Join("page", pageName),
			Name:     fmt.Sprintf("handlers_%X.go", time.Now().Unix()),
			Contents: mustExecute("scaffold/page/name/handlers.go.tmpl", data, scaffold),
		})
	}
	for i := 0; i < len(templates); i++ {
		tmpl := templates[i]
		files = append(files, File{
			Dir:      filepath.Dir(tmpl.Path),
			Name:     filepath.Base(tmpl.Path),
			Contents: mustExecute("scaffold/page/name/templates/block/partial.go.tmpl", data, scaffold),
		})
	}
	return
}
