package generate

import (
	"io/fs"
	"path/filepath"
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

// RoutesScaffold will generate all files for the good routes command
func RoutesScaffold(sitePkg GoPackage, pageName string, entries []Entry, routes []Route, scaffold fs.FS) (files []File, err error) {
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
	return
}
