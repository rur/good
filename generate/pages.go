package generate

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
)

// PagesFile creates a new pages.go file by scanning the target scaffold site
// for directories inside the ./page sub-package
func PagesFile(sitePkg GoPackage, scaffold fs.FS) (file File, err error) {
	list, err := ioutil.ReadDir(path.Join(sitePkg.Dir, "page"))
	if err != nil {
		err = fmt.Errorf("failed to scan scaffold '%s' for pages: %s", sitePkg.ImportPath, err)
		return
	}
	var pages []string
	for i := range list {
		// note that 'templates' is reserved for shared template files
		if list[i].IsDir() && list[i].Name() != "templates" {
			pages = append(pages, list[i].Name())
		}
	}
	// treat each dir name as a page
	file.Dir, err = filepath.Rel(sitePkg.Module.Dir, sitePkg.Dir)
	if err != nil {
		return
	}
	file.Name = "pages.go"
	file.Contents = mustExecute("scaffold/pages.go.tmpl", struct {
		Pages     []string
		Namespace string
	}{
		Pages:     pages,
		Namespace: sitePkg.ImportPath,
	}, scaffold)
	return
}
