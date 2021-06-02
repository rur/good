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

// ScaffoldPage will assemble files for adding a new page to the site scaffold
func ScaffoldPage(sitePkg GoPackage, name string, scaffold fs.FS) (files []File, err error) {
	relDir, err := sitePkg.RelPath()
	if err != nil {
		return
	}

	// setup page with some placeholder data
	data := struct {
		Name       string // Go package name for page
		Namespace  string
		SiteDirRel string
		Handlers   []Handler
		PageEntry  Entry
		Entries    []Entry
		Routes     []Route
		Templates  string
		PagePath   string
	}{
		PagePath:   strings.Join([]string{sitePkg.ImportPath, "page", name}, "/"),
		Name:       name,
		Namespace:  sitePkg.ImportPath,
		SiteDirRel: relDir,
		Handlers: []Handler{
			{
				Ref:        name,
				Type:       "DefaultSubView",
				Extends:    "content",
				Method:     "GET",
				Doc:        "Root handler for the " + name + " page",
				Identifier: name + "Handler",
				Blocks: []HandleBlock{
					{FieldName: "SiteNav", Name: "site-nav"},
					{FieldName: "Content", Name: "content"},
					{FieldName: "Scripts", Name: "scripts"},
				},
			},
			{
				Ref:        "placeholder",
				Type:       "DefaultSubView",
				Extends:    "content",
				Method:     "GET",
				Doc:        "This is placeholder content, add your endpoints to the routemap.toml and run go generate",
				Identifier: "placeholderHandler",
			},
		},
		PageEntry: Entry{
			Assignment: name,
			Template:   filepath.Join("page", name, "templates", name+".html.tmpl"),
			Handler:    fmt.Sprintf("hlp.BindEnv(bindResources(%sHandler))", name),
		},
		Entries: []Entry{{
			Assignment: "",
			Block:      "site-nav",
			Type:       "DefaultSubView",
			Extends:    name,
			Template:   filepath.Join("page", "templates", "nav.html.tmpl"),
			Handler:    "hlp.BindEnv(page.SiteNavHandler)",
		}, {
			Assignment: "placeholder",
			Block:      "content",
			Type:       "DefaultSubView",
			Extends:    name,
			Template:   filepath.Join("page", name, "templates", "content", "placeholder.html.tmpl"),
			Handler:    "hlp.BindEnv(bindResources(placeholderHandler))",
		}, {
			Assignment: "",
			Block:      "scripts",
			Type:       "DefaultSubView",
			Extends:    name,
			Template:   filepath.Join("page", "templates", "scripts.html.tmpl"),
			Handler:    "treetop.Noop",
		}},
		Routes: []Route{{
			Method:    "GET",
			Path:      "/" + name,
			Reference: "placeholder",
		}},
		Templates: filepath.Join("page", name, "templates"),
	}

	pageDir := filepath.Join(relDir, "page", name)

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
	// page/name/routes.go
	files = append(files, File{
		Dir:      pageDir,
		Name:     "routes.go",
		Contents: mustExecute("scaffold/page/name/routes.go.tmpl", data, scaffold),
	})
	// page/name/templates/name.html.tmpl
	files = append(files, File{
		Dir:      filepath.Join(pageDir, "templates"),
		Name:     name + ".html.tmpl",
		Contents: mustExecute("scaffold/page/name/templates/name.html.tmpl.tmpl", data, scaffold),
	})
	// page/name/templates/content/placeholder.html.tmpl
	files = append(files, File{
		Dir:      filepath.Join(pageDir, "templates", "content"),
		Name:     "placeholder.html.tmpl",
		Contents: mustExecute("scaffold/page/name/templates/content/placeholder.html.tmpl.tmpl", data, scaffold),
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
