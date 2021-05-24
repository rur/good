package generate

import (
	"embed"
	"io/fs"
	"path"
	"testing"
)

//go:embed testdata/scaffold
var testScaffold embed.FS

func TestSiteScaffold(t *testing.T) {
	scaffold, _ := fs.Sub(testScaffold, "testdata")

	pkg, err := GoListPackage(".")
	if err != nil {
		t.Fatal("ValidateScaffoldPackage error getting local package details", err)
	}
	sitePkg, siteDir, err := ValidateScaffoldPackage(
		pkg.Module, "./admin/site", scaffold,
	)

	if err != nil {
		t.Error("ValidateScaffoldPackage: unexpected error", err)
		return
	}
	expect := "github.com/rur/good/admin/site"
	if sitePkg != expect {
		t.Errorf("ValidateScaffoldPackage: expecting site package %s, got %s", expect, sitePkg)
	}
	expect = path.Join(pkg.Module.Dir, "admin", "site")
	if siteDir != expect {
		t.Errorf("ValidateScaffoldPackage: expecting site directory %s, got %s", expect, siteDir)
	}
}
