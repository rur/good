package generate

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSiteFromPagePackage(t *testing.T) {
	tests := []struct {
		name        string
		pkg         GoPackage
		wantSitePkg GoPackage
		wantErr     bool
	}{
		{
			name: "typical",
			pkg: GoPackage{
				Dir:        "/some/path/to/module/site/page/test",
				ImportPath: "github.com/rur/example/site/page/test",
				Module: GoModule{
					Path: "github.com/rur/example",
					Dir:  "/some/path/to/module",
				},
			},
			wantSitePkg: GoPackage{
				Dir:        "/some/path/to/module/site",
				ImportPath: "github.com/rur/example/site",
				Module: GoModule{
					Path: "github.com/rur/example",
					Dir:  "/some/path/to/module",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSitePkg, err := SiteFromPagePackage(tt.pkg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SiteFromPagePackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotSitePkg, tt.wantSitePkg) {
				t.Errorf("SiteFromPagePackage() = %v, want %v", gotSitePkg, tt.wantSitePkg)
			}
		})
	}
}

func TestPageScaffold(t *testing.T) {
	fs := os.DirFS("../")
	gotFiles, err := PageScaffold(
		GoPackage{
			ImportPath: "github.com/rur/example/admin/site",
			Dir:        "/some/location/example/admin/site",
			Module: GoModule{
				Dir:  "/some/location/example",
				Path: "github.com/rur/example",
			},
		},
		"testing",
		fs,
	)
	if err != nil {
		t.Errorf("PageScaffold() error = %v", err)
		return
	}
	gotFileMap := filesToMap(gotFiles)

	expectedFileMap := map[string][]string{
		"page/testing/resources.go": {
			`"github.com/rur/example/admin/site/page"`,
		},
		"page/testing/routemap.toml": {
			`_template = "page/testing/templates/testing.html.tmpl`,
			`_handler = "hlp.BindEnv(bindResources(testingHandler))"`,
		},
		"page/testing/templates/testing.html.tmpl": {
			`{{ template "site-nav" .SiteNav }}`,
			`{{ template "content" .Content }}`,
			`{{ template "scripts" .Scripts }}`,
		},
		"page/testing/templates/content/placeholder.html.tmpl": {
			"<h1>Run go generate command for page testing</h1>",
		},
		"page/testing/gen.go": {
			"//go:generate good routes .",
		},
		"page/testing/handlers.go": {
			`SiteNav:`,
			`rsp.HandleSubView("site-nav", req)`,
		},
	}

	for file, checks := range expectedFileMap {
		content, ok := gotFileMap[file]
		if !ok {
			t.Errorf("Expecting file '%s'", file)
		} else {
			for _, pattern := range checks {
				if !strings.Contains(content, pattern) {
					t.Errorf("Expecting file '%s' to contain %s\nGOT: %s", file, pattern, content)
				}
			}
		}
	}
	for file := range gotFileMap {
		_, ok := expectedFileMap[file]
		if !ok {
			t.Errorf("Unexpected file '%s'", file)
		}
	}
}
