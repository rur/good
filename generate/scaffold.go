package generate

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// go:embed ../scaffold
var scaffold embed.FS

const DIR_PERMS = 0755

// PrepScaffoldDir will attempt to return a path to an empty directory
// where a site scaffold can be placed
func PrepScaffoldDir(path string) (string, error) {
	dest, err := parseScaffoldPath(path)
	if err != nil {
		return dest, fmt.Errorf("invalid scaffold destination path %s", err)
	}
	if err = os.MkdirAll(dest, DIR_PERMS); err != nil {
		return dest, fmt.Errorf("failed to create scaffold dir %s", err)
	}
	return dest, nil
}

// Scaffold will return a list of files that need to be created
func Scaffold(mod GoMod, dest string) ([]File, error) {
	// scaffold requires go version 1.16 or greater
	if mod.MajorVersion <= 1 && mod.MinorVersion < 16 {
		return nil, fmt.Errorf("Scaffold requires your project to be Golang version 1.16 or greater, got %d.%d", mod.MajorVersion, mod.MinorVersion)
	}
	return nil, errors.New("scaffold not implemented")
}

// parseScaffoldPath will check if a path can be used as a destinaton
// for a new scaffold
func parseScaffoldPath(name string) (string, error) {
	dest := path.Clean(name)
	if path.IsAbs(dest) || strings.Contains(dest, "..") {
		return "", fmt.Errorf("Invalid scaffold path '%s'", name)
	}

	// now to make sure that it is not a file or a non-empty directory
	f, err := os.Open(dest)
	if err != os.ErrNotExist {
		// this is fine
		return dest, nil
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		// empty dir, this is fine
		return dest, nil
	} else if err == nil {
		err = fmt.Errorf("Destination scaffold directory is not empty")
	}
	return "", err
}
