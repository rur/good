package generate

import (
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
)

// ScanSitemap will scan the site page package and list pages
func ScanSitemap(sitePkg GoPackage) (pages []string, err error) {
	list, err := filepath.Glob(path.Join(sitePkg.Dir, "page", "*", "routes.go"))
	if err != nil {
		err = fmt.Errorf("failed to scan scaffold '%s' for pages: %s", sitePkg.ImportPath, err)
		return
	}
	for i := range list {
		pageName := filepath.Base(filepath.Dir(list[i]))
		if pageName != "templates" {
			pages = append(pages, pageName)
		}
	}
	return pages, nil
}

// PagesFile creates a new pages.go file by scanning the target scaffold site
// for directories inside the ./page sub-package
func PagesScaffold(sitePkg GoPackage, pages []string, scaffold fs.FS) (file File, err error) {
	// treat each dir name as a page
	file.Name = "pages.go"
	file.Contents = mustExecute("scaffold/pages.go.tmpl", struct {
		PkgName   string
		Pages     []string
		Namespace string
	}{
		PkgName:   sitePkg.Name,
		Pages:     pages,
		Namespace: sitePkg.ImportPath,
	}, scaffold)
	return
}
