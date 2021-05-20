package generate

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

// ModFile some info parsed from a go.mod file
type GoMod struct {
	Module       string
	Version      string
	MajorVersion int
	MinorVersion int
}

// ReadGoModFile will read the go.mod file from the current working directory or panic
func ReadGoModFile(path string) (GoMod, error) {
	var mod GoMod
	f, err := os.Open(path)
	if err != nil {
		return mod, err
	}
	scan := bufio.NewReader(f)
	for {
		line, err := scan.ReadString('\n')
		if err == io.EOF {
			// end of go mod file, return what we have parsed so far
			return mod, nil
		} else {
			if err != nil {
				return mod, err
			}
		}
		if strings.HasPrefix(line, "module ") {
			mod.Module = strings.TrimSpace(strings.Split(line, " ")[1])
		}
		if strings.HasPrefix(line, "go ") {
			mod.Version = strings.TrimSpace(strings.Split(line, " ")[1])
			parts := strings.Split(mod.Version, ".")
			mod.MajorVersion, err = strconv.Atoi(parts[0])
			if err != nil {
				return mod, err
			}
			mod.MinorVersion, err = strconv.Atoi(parts[1])
			if err != nil {
				return mod, err
			}
		}
	}
}
