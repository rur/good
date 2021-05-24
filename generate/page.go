package generate

import (
	"io/fs"
	"path/filepath"
)

// ScaffoldPage will assemble files for adding a new page to the site scaffold
func ScaffoldPage(siteModule, siteDir, name string, scaffold fs.FS) (files []File, err error) {
	data := struct {
		Namespace string
	}{
		Namespace: siteModule,
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
	// page/name/routemap.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "routemap.go",
		Contents: mustExecute("scaffold/page/name/routemap.go.tmpl", data, scaffold),
	})
	// page/name/templates/placeholder.html.tmpl
	files = append(files, File{
		Dir:      filepath.Join(pageDir, "templates"),
		Name:     "placeholder.html.tmpl",
		Contents: mustExecute("scaffold/page/name/templates/placeholder.html.tmpl", data, scaffold),
	})
	return
}
