package generate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/modfile"
)

// duration that calls to exec.Command are allow to block before they are cancelled
const execTimeout = 3 * time.Second

// GoModule represents a golang project level module
type GoModule struct {
	Path string `json:"Path"`
	Dir  string `json:"Dir"`
}

// GoPackage represents a package within a go module
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

// IsTTY will check if we have a terminal, should work on Linux or Mac
func IsTTY() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// Name will return the name of the package based on the import path
func (pkg *GoPackage) Name() string {
	ind := strings.LastIndex(pkg.ImportPath, "/")
	return pkg.ImportPath[ind+1:]
}

// GoListPackage will get the Go module information for the go path provied
func GoListPackage(path string) (pkg GoPackage, err error) {
	var stdout, stderr bytes.Buffer
	timeout, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()
	if path == "./..." {
		cmd := exec.CommandContext(timeout, "go", "list", "./...")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err = cmd.Run(); err != nil {
			err = fmt.Errorf("%s: %s", err, stderr.String())
			return
		}
		firstPath, stdOutErr := stdout.ReadString('\n')
		if stdOutErr == nil || stdOutErr == io.EOF {
			path = strings.TrimSpace(firstPath)
		} else {
			err = fmt.Errorf("failed to find a valid golang module in this directory, got output: %s, with error: %s", firstPath, stdOutErr)
			return
		}
	}

	cmd := exec.CommandContext(timeout, "go", "list", "-e", "--json", path)
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%s: %s", err, stderr.String())
		return
	}
	target := struct {
		GoPackage
		Error struct {
			Err string
		}
	}{}
	err = json.Unmarshal(stdout.Bytes(), &target)
	if target.Dir == "" {
		err = fmt.Errorf("failed to load a go package for path '%s'", path)
		return
	} else {
		pkg = target.GoPackage
	}
	if err != nil || pkg.Module.Dir != "" {
		// we have everything we need
		return
	}

	// failed to load the module, the target package is probably malformed
	// try to load the module by trimming the package import path
	parts := strings.Split(pkg.ImportPath, "/")
	for i := len(parts); i > 0; i-- {
		pkg.Module, err = GoListModule(strings.Join(parts[:i], "/"))
		if err == nil {
			return
		}
	}
	err = fmt.Errorf("failed to load package module for path %s", path)
	return
}

// GoListModule will get the Go module information for the import path provied
func GoListModule(path string) (mod GoModule, err error) {
	var stdout, stderr bytes.Buffer
	timeout, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()
	cmd := exec.CommandContext(timeout, "go", "list", "-m", "--json", path)
	stdout.Reset()
	stderr.Reset()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%s: %s", err, stderr.String())
		return
	}
	err = json.Unmarshal(stdout.Bytes(), &mod)
	return
}

// GoFormat will execute the go fmt command on the module path
func GoFormat(path string) (string, error) {
	timeout, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()
	cmd := exec.CommandContext(timeout, "go", "fmt", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go fmt error: %s, output: %s", err, output)
	}
	return string(output), nil
}

// TryReadModFile will attempt to load the root module path from a go.mod file
// it will return empty string otherwise
func TryReadModFile() (module string, baseDir string) {
	defer func() {
		if module == "" {
			baseDir = ""
		}
	}()
	var err error
	baseDir, err = os.Getwd()
	if err != nil {
		return
	}
	contents, err := ioutil.ReadFile(filepath.Join(baseDir, "go.mod"))
	if err != nil {
		return
	}
	module = modfile.ModulePath(contents)
	return
}

// TimeoutScanln will wait for at most 10 seconds for a line of input from the user
func TimeoutScanln() (string, error) {
	var (
		input string
		err   error
		done  = make(chan struct{})
	)
	go func() {
		_, err = fmt.Scanln(&input)
		close(done)
	}()
	select {
	case <-done:
		return input, err
	case <-time.After(30 * time.Second):
		return "", errors.New(
			"input timeout, sorry you'll have to be quicker than that! " +
				"Good CLI will not block the terminal for more than 30 seconds",
		)
	}
}
