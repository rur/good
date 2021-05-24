package generate

import (
	"os"
	"path/filepath"
)

const FILE_PERMS = 0755

// File description and contents that will later be
// flushed to disk
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

// flushFile will attempt to build the necessary folder
// structure and write content bytes to disk
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
