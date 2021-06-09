package generate

import (
	"os"
	"path/filepath"
)

const (
	FOLDER_PERMS = 0755
	FILE_PERMS   = 0644
)

// File description and contents that will later be
// flushed to disk
type File struct {
	Dir       string
	Name      string
	Contents  []byte
	Overwrite bool
}

// Path will return the dir+name for this file
func (f File) Path() string {
	return filepath.Join(f.Dir, f.Name)
}

// FlushFiles will write files out to a temp directory and move those files
// to into the go package
func FlushFiles(modDir string, files []File) error {
	for _, file := range files {
		err := flushFile(modDir, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// flushFile will attempt to build the necessary folder
// structure and write content bytes to disk
func flushFile(dir string, file File) error {
	err := os.MkdirAll(filepath.Join(dir, file.Dir), FOLDER_PERMS)
	if err != nil {
		return err
	}
	mode := os.O_RDWR | os.O_CREATE
	if file.Overwrite {
		mode = mode | os.O_TRUNC
	}
	fh, err := os.OpenFile(filepath.Join(dir, file.Path()), mode, FILE_PERMS)
	if err != nil {
		return err
	}
	defer fh.Close()
	fh.Truncate(0)
	fh.Seek(0, 0)
	_, err = fh.Write(file.Contents)
	return err
}
