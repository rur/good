package generate

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestGoListPackageGenerate(t *testing.T) {
	// helpers
	assertPattern := func(field, got, expected string) bool {
		if !regexp.MustCompile(expected).MatchString(got) {
			t.Errorf("GoListPackage pkg.%s, expected pattern %s, got %s", field, expected, got)
			return false
		}
		return true
	}
	assertPath := func(field, got, pattern string) bool {
		return assertPattern(
			field,
			got,
			strings.ReplaceAll(pattern, `/`, string(os.PathSeparator)),
		)
	}

	pkg, err := GoListPackage(".")

	if err != nil {
		t.Errorf("GoListPackage() error %s", err)
		return
	}
	if !assertPath("Dir", pkg.Dir, `/good/generate$`) {
		return
	}
	if !assertPattern("ImportPath", pkg.ImportPath, "github.com/rur/good/generate") {
		return
	}
	if !assertPath("Root", pkg.Root, `/good$`) {
		return
	}
	if !assertPattern("Module.Path", pkg.Module.Path, "^github.com/rur/good$") {
		return
	}
	if !assertPath("Module.Dir", pkg.Module.Dir, "/good$") {
		return
	}
	if !assertPath("Module.GoMod", pkg.Module.GoMod, `go\.mod$`) {
		return
	}
	if !assertPattern("Module.GoVersion", pkg.Module.GoVersion, `\d.\d{2}`) {
		return
	}
}

func TestGoListPackageError(t *testing.T) {
	// fs error
	if _, err := GoListPackage("./fake"); err == nil {
		t.Error("Expecting error for fake directory")
	} else if !strings.Contains(err.Error(), "directory not found") {
		t.Errorf("Expecting reason for failure, got %s", err)
	}
	// invoke a go command error
	if _, err := GoListPackage("__does_not_exist_in_go_root__"); err == nil {
		t.Error("Expecting error for fake directory")
	} else if !strings.Contains(err.Error(), "package __does_not_exist_in_go_root__ is not in GOROOT") {
		t.Errorf("Expecting reason for failure, got %s", err)
	}
}
