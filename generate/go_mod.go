package generate

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// ModFile some info parsed from a go.mod file
type ModFile struct {
	Module  string
	Version string
}

// MustReadGoModFile will read the go.mod file from the current working directory or panic
func MustReadGoModFile(path string) ModFile {
	var mod ModFile
	f, err := os.Open(path)
	mustNot(err)
	scan := bufio.NewReader(f)
	for {
		line, err := scan.ReadString('\n')
		if err == io.EOF {
			// end of go mod file, return what we have parsed so far
			return mod
		} else {
			mustNot(err)
		}
		if strings.HasPrefix(line, "module ") {
			mod.Module = strings.TrimSpace(strings.Split(line, " ")[1])
		}
		if strings.HasPrefix(line, "go ") {
			mod.Version = strings.TrimSpace(strings.Split(line, " ")[1])
		}
	}
}

// mustNot will panic of err value is not nil
func mustNot(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
