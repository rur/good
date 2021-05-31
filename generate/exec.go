package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type GoModule struct {
	Path      string `json:"Path"`
	Main      bool   `json:"Main"`
	Dir       string `json:"Dir"`
	GoMod     string `json:"GoMod"`
	GoVersion string `json:"GoVersion"`
}

type GoPackage struct {
	Dir        string   `json:"Dir"`
	ImportPath string   `json:"ImportPath"`
	Name       string   `json:"Name"`
	Root       string   `json:"Root"`
	Module     GoModule `json:"Module"`
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

// GoFormat will execute the go fmt command on the module path
func GoFormat(path string) (string, error) {
	cmd := exec.Command("go", "fmt", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go fmt error: %s, output: %s", err, output)
	}
	return string(output), nil
}

// GoGenerate will execute the go generate command on the module path
func GoGenerate(path string) (string, error) {
	cmd := exec.Command("go", "generate", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go fmt error: %s, output: %s", err, output)
	}
	return string(output), nil
}
