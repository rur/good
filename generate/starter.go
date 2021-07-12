package generate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// StarterScaffold will output a
func StarterScaffold(dir string, scaffold fs.FS) (files []File, err error) {
	_, err = os.Stat(dir)
	if !os.IsNotExist(err) {
		if err == nil {
			return nil, fmt.Errorf("File path '%s' already exists, please provide a valid location for a new folder", dir)
		}
		return
	} else {
		err = nil
	}

	err = fs.WalkDir(scaffold, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		contents, err := fs.ReadFile(scaffold, path)
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

	return

}
