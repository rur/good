package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// GoModule represents a golang project level module
type GoModule struct {
	Path  string `json:"Path"`
	Dir   string `json:"Dir"`
	GoMod string `json:"GoMod"`
}

// GoPackage represents a package with a go module
type GoPackage struct {
	Dir        string   `json:"Dir"`
	ImportPath string   `json:"ImportPath"`
	Module     GoModule `json:"Module"`
}

// RelPath gets the relative import path for the package relative
// to the module, if will return "." if they are the same
func (pkg *GoPackage) RelPath() (string, error) {
	if pkg.Dir == pkg.Module.Dir {
		return ".", nil
	}
	path, err := filepath.Rel(pkg.Module.Dir, pkg.Dir)
	return "./" + path, err
}

// Name will return the name of the package based on the import path
func (pkg *GoPackage) Name() string {
	ind := strings.LastIndex(pkg.ImportPath, "/")
	return pkg.ImportPath[ind+1:]
}

// GoListPackage will get the Go module information for the go path provied
func GoListPackage(path string) (pkg GoPackage, err error) {
	if path == "" {
		path = "."
	} else if path[0] != '.' {
		path = "./" + path
	}

	cmd := exec.Command("go", "list", "--json", path)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%s: %s", err, stderr.String())
		return
	}
	err = json.Unmarshal(stdout.Bytes(), &pkg)
	return
}

// GoFormat will execute the go fmt command on the module path
func GoFormat(path string) (string, error) {
	cmd := exec.Command("go", "fmt", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go fmt error: %s, output: %s", err, output)
	}
	return string(output), nil
}
