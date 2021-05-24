package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type GoModule struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

type GoPackage struct {
	Dir        string
	ImportPath string
	Name       string
	Root       string
	Module     GoModule
}

// GoListPackage will get the Go module information for the go path provied
func GoListPackage(path string) (pkg GoPackage, err error) {
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
