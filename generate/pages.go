package generate

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
)

// ScanSitemap will scan the site page package and load site routing
// data
func ScanSitemap(sitePkg GoPackage) (pages []string, err error) {
	dir := path.Join(sitePkg.Dir, "page")
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		err = fmt.Errorf("failed to scan scaffold '%s' for pages: %s", sitePkg.ImportPath, err)
		return
	}
	for i := range list {
		// note that 'templates' is reserved for shared template files
		name := list[i].Name()
		if list[i].IsDir() && name != "templates" {
			pages = append(pages, list[i].Name())
		}
	}
	return pages, nil
}

// PagesFile creates a new pages.go file by scanning the target scaffold site
// for directories inside the ./page sub-package
func PagesScaffold(sitePkg GoPackage, pages []string, scaffold fs.FS) (file File, err error) {
	// treat each dir name as a page
	file.Dir, err = filepath.Rel(sitePkg.Module.Dir, sitePkg.Dir)
	if err != nil {
		return
	}
	file.Name = "pages.go"
	file.Contents = mustExecute("scaffold/pages.go.tmpl", struct {
		Pages       []string
		Namespace   string
		SitemapJSON string
	}{
		Pages:     pages,
		Namespace: sitePkg.ImportPath,
		SitemapJSON: `{
			"example": {
				"path": "/example",
				"routes": {
					"placeholder": {
						"block": "content",
						"path": "/example/placeholder"
					}
				}
			},
			"testing: {
				"path": "/testing",
				"routes": {
					"placeholder": {
						"block": "content",
						"path": "/testing/placeholder"
					}
				}
			}
		}`,
	}, scaffold)
	return
}
