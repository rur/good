package generate

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	PageNameRegex = regexp.MustCompile(`^[a-z][a-z]+$`)
)

// PageScaffold will assemble files for adding a new page to the site scaffold
func PageScaffold(sitePkg GoPackage, name string, scaffold fs.FS) (files []File, err error) {
	// setup page with some placeholder data
	data := struct {
		Name      string // Go package name for page
		Namespace string
		Handlers  []Handler
		Templates string
		PagePath  string
	}{
		PagePath:  strings.Join([]string{sitePkg.ImportPath, "page", name}, "/"),
		Name:      name,
		Namespace: sitePkg.ImportPath,
		Handlers: []Handler{
			{
				Ref:        "example-dummy",
				Block:      "content",
				Method:     "GET",
				Doc:        "This is an unused handler for the sake of example",
				Identifier: "exampleDummyHandler",
			},
		},
		Templates: filepath.Join("page", name, "templates"),
	}

	pageDir := filepath.Join("page", name)

	// page/name/gen.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "gen.go",
		Contents: mustExecute("scaffold/page/name/gen.go.tmpl", data, scaffold),
	})
	// page/name/handlers.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "handlers.go",
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
	// page/name/templates/placeholder.html.tmpl
	files = append(files, File{
		Dir:      filepath.Join(pageDir, "templates"),
		Name:     "placeholder.html.tmpl",
		Contents: mustExecute("scaffold/page/name/templates/placeholder.html.tmpl.tmpl", data, scaffold),
	})

	return
}

// ValidatePagePackage will check that a page can be installed at a given location.
//
// If the target directory is not empty, this will scan for conflicts against the scaffold
func ValidatePageLocation(pageDir string, scaffold fs.FS) error {
	// now try to check if there will be files write conflicts
	// build block list index
	blocked := struct{}{} // zero size sentinel
	blockList := make(map[string]struct{})
	entries, err := fs.ReadDir(scaffold, filepath.Join("scaffold", "page", "name"))
	if err != nil {
		return fmt.Errorf("failed to open scaffold/page/name folder: %s", err)
	}
	for _, entry := range entries {
		blockList[strings.TrimSuffix(entry.Name(), ".tmpl")] = blocked
	}

	// Scan for conflict between the scaffold and the target FS
	fh, err := os.Open(pageDir)
	if os.IsNotExist(err) {
		// no destination folder, no conflicts, all good
		return nil
	} else if err != nil {
		return fmt.Errorf("error while scanning target dir: %s", err)
	}
	defer fh.Close()
	// Limit scan to 500 folder children
	names, err := fh.Readdirnames(500)
	if err == io.EOF {
		// empty dir, no conflicts, all good
		return nil
	}
	for i := 0; i < len(names); i++ {
		if _, ok := blockList[names[i]]; ok {
			return fmt.Errorf("conflicting file or direcotry '%s'", names[i])
		}
	}
	return nil
}

// ValidatePageName will check if a page name is permitted to be the name of a site page
func ValidatePageName(name string) error {
	if name == "templates" {
		return errors.New("'templates' cannot be used as a page name, it is reserved for the shared template directory")
	}
	if !PageNameRegex.MatchString(name) {
		return fmt.Errorf(
			`page name '%s' is not valid. Use best practices for Go package names: all lowercase, all alpha. See helpful guideline https://blog.golang.org/package-names`,
			name,
		)
	}
	return nil
}

// SiteFromPagePackage will construct a site go package given one that refers to a page in that
// site
func SiteFromPagePackage(pkg GoPackage) (sitePkg GoPackage, err error) {
	parts := strings.Split(pkg.Dir, string(os.PathSeparator))
	for parts[len(parts)-1] == "" {
		// trim trailing slashes
		parts = parts[:len(parts)-1]
	}
	if len(parts) < 3 || parts[len(parts)-2] != "page" {
		err = fmt.Errorf("unexpected page package path: '%s'", pkg.ImportPath)
		return
	}
	// drop '/page/name' suffix to get site path
	dir := strings.Join(parts[:len(parts)-2], string(os.PathSeparator))
	sitePkg = GoPackage{
		Dir:        dir,
		ImportPath: strings.TrimSuffix(pkg.ImportPath, "/page/"+parts[len(parts)-1]),
		Module:     pkg.Module,
	}
	return
}
