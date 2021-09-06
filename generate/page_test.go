package generate

import (
	"os"
	"path/filepath"
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
	scaffold := os.DirFS("..")
	bootstrap := os.DirFS(filepath.Join("testdata", "bootstrap"))

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
		scaffold,
		bootstrap,
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
			`_handler = 'hlp.BindEnv(page.GetBaseHandler("testing Page"))'`,
			`_template = "page/templates/scripts.html.tmpl"`,
		},
		"page/testing/templates/placeholder.html.tmpl": {
			`$ go generate github.com/rur/example/admin/site/...`,
		},
		"page/testing/gen.go": {
			"//go:generate good routes gen .",
		},
		"page/testing/routes.go": {
			"hlp.BindEnv(bindResources(readmePageHandler)),",
		},
		"page/testing/README.md": {
			"# Page `testing`",
		},
		"page/testing/handlers.go": {
			`"github.com/rur/example/admin/site/service"`,
			`(rsc *resources, env *service.Env, rsp treetop.Response, req *http.Request) interface{} {`,
		},
		"static/js/testing/example.js": {
			`function test() {`,
		},
		"static/styles/testing/example.css": {
			`.test {`,
		},
		"static/public/testing/example.txt": {
			`This is a test`,
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
