package generate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// StarterScaffold will output a copy of the raw starter template files
func StarterScaffold(dir string, scaffold, starter fs.FS) (files []File, err error) {
	_, err = os.Stat(dir)
	if !os.IsNotExist(err) {
		if err == nil {
			return nil, fmt.Errorf("File path '%s' already exists, please provide a valid location for a new folder", dir)
		}
		return
	} else {
		err = nil
	}
	var found = map[string]bool{
		"routes.go.tmpl":     false,
		"routemap.toml.tmpl": false,
		"resources.go.tmpl":  false,
		"handlers.go.tmpl":   false,
		"README.md.tmpl":     false,
	}
	err = fs.WalkDir(starter, ".", func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		found[path] = true

		contents, err := fs.ReadFile(starter, path)
		if err != nil {
			return err
		}

		files = append(files, File{
			Dir:      filepath.Join(dir, filepath.Dir(path)),
			Name:     filepath.Base(path),
			Contents: contents,
		})
		return nil
	})

	for path, found := range found {
		if !found {
			var contents []byte
			contents, err = fs.ReadFile(scaffold, filepath.Join("scaffold", "page", "default", path))
			if err != nil {
				return
			}
			files = append(files, File{
				Dir:      filepath.Join(dir, filepath.Dir(path)),
				Name:     filepath.Base(path),
				Contents: contents,
			})
		}
	}

	return

}
