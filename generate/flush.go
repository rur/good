package generate

import (
	"fmt"
	"io/ioutil"
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

// FlushFiles will write files out to a temp directory and move those files
// to into the go package
func FlushFiles(modulePath string, files []File) error {
	tmp, err := ioutil.TempDir("", "goob-generate-")
	if err != nil {
		return fmt.Errorf("failed to create a temporary directory: %s", err)
	}
	for _, file := range files {
		err := flushFile(tmp, file)
		if err != nil {
			return err
		}
	}
	items, _ := ioutil.ReadDir(tmp)
	for _, item := range items {
		err = os.Rename(filepath.Join(tmp, item.Name()), filepath.Join(modulePath, item.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

// flushFile will attempt to build the necessary folder
// structure and write content bytes to disk
func flushFile(tmp string, file File) error {
	err := os.MkdirAll(filepath.Join(tmp, file.Dir), FILE_PERMS)
	if err != nil {
		return err
	}
	fh, err := os.Create(filepath.Join(tmp, file.Path()))
	if err != nil {
		return err
	}
	defer fh.Close()
	_, err = fh.Write(file.Contents)
	return err
}
