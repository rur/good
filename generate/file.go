package generate

import "fmt"

type File struct {
	RelBaseDir []string
	Name       string
	Ext        string
	Contents   []byte
}

// FlushFiles will write files to FS
func FlushFiles(files []File) error {
	fmt.Printf("Write %d file", len(files))
	return nil
}
