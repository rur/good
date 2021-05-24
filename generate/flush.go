package generate

import (
	"os"
	"path/filepath"
)

const FILE_PERMS = 0755

type File struct {
	Dir      string
	Name     string
	Contents []byte
}

// Path will return the dir+name for this file
func (f File) Path() string {
	return filepath.Join(f.Dir, f.Name)
}

// FlushFiles will write files out to the file system
func FlushFiles(files []File) error {
	for _, file := range files {
		err := flushFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func flushFile(file File) error {
	err := os.MkdirAll(file.Dir, FILE_PERMS)
	if err != nil {
		return err
	}
	fh, err := os.Create(file.Path())
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = fh.Write(file.Contents)
	return err
}
