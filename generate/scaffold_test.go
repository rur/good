package generate

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

//go:embed testdata/scaffold
var testScaffold embed.FS

func TestValidateScaffoldPackage(t *testing.T) {
	pkg, err := GoListPackage(".")
	if err != nil {
		t.Fatal(err)
	}
	scaffold, _ := fs.Sub(testScaffold, "testdata")
	if err != nil {
		t.Fatal("ValidateScaffoldPackage error getting local package details", err)
	}
	tests := []struct {
		name     string
		location string
		wantErr  string
	}{
		{
			name:     "basic",
			location: "./admin/site",
		},
		{
			name:     "using . as destination",
			location: ".",
		},
		{
			name:     "conflicting file",
			location: "./generate/testdata/with_conflict_file",
			wantErr:  "conflicting file or directory 'file.go'",
		},
		{
			name:     "conflicting directory",
			location: "./generate/testdata/with_conflict_folder",
			wantErr:  "conflicting file or directory 'folder'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScaffoldLocation(filepath.Join(pkg.Module.Dir, tt.location), scaffold)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("ValidateScaffoldPackage() expecting an error containing message %s", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("ValidateScaffoldPackage() expecting error to contain '%s', got '%s'", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("ValidateScaffoldPackage() unexpected error = %s", err)
			}
		})
	}
}

func TestParseSitePackage(t *testing.T) {
	pkg, err := GoListPackage(".")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		input       string
		wantPkg     string
		wantDir     string
		wantRelPath string
		wantErr     string
	}{
		{
			name:        "basic",
			input:       "./admin/site",
			wantPkg:     "github.com/rur/good/admin/site",
			wantDir:     filepath.Join(pkg.Module.Dir, "admin", "site"),
			wantRelPath: "./admin/site",
		},
		{
			name:        "using . as destination",
			input:       ".",
			wantPkg:     "github.com/rur/good",
			wantDir:     pkg.Module.Dir,
			wantRelPath: ".",
		},
		{
			name:    "embedded import",
			input:   "github.com/rur/good/admin/site",
			wantErr: "site package name must be relative to the current module, got github.com/rur/good/admin/site",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPkg, err := ParseSitePackage(pkg.Module, tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("ParseSitePackage() expecting an error containing message %s", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("ParseSitePackage() expecting error to contain '%s', got '%s'", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("ParseSitePackage() unexpected error = %s", err)
			} else {
				if gotPkg.ImportPath != tt.wantPkg {
					t.Errorf("ParseSitePackage() gotSitePkg = %v, want %v", gotPkg.ImportPath, tt.wantPkg)
				}
				if gotPkg.Dir != tt.wantDir {
					t.Errorf("ParseSitePackage() gotSiteDir = %v, want %v", gotPkg.Dir, tt.wantDir)
				}
				rel, _ := gotPkg.RelPath()
				if rel != tt.wantRelPath {
					t.Errorf("ParseSitePackage() RelPath() = %v, want %v", rel, tt.wantRelPath)
				}
			}
		})
	}
}

func TestSiteScaffold(t *testing.T) {
	fs := os.DirFS("../")
	gotFiles, err := SiteScaffold(
		GoPackage{
			ImportPath: "github.com/rur/example/admin/site",
			Dir:        "/some/location/example/admin/site",
			Module: GoModule{
				Dir:  "/some/location/example",
				Path: "github.com/rur/example",
			},
		},
		fs,
	)
	if err != nil {
		t.Errorf("SiteScaffold() error = %v", err)
		return
	}
	gotFileMap := filesToMap(gotFiles)
	tests := []struct {
		name     string
		file     string
		contains string
	}{
		{
			"simple /site/user.go",
			"site/user.go",
			"type User struct {",
		}, {
			"simple /page/helper.go",
			"page/helper.go",
			`"github.com/rur/example/admin/site/site"`,
		}, {
			"simple /main.go",
			"main.go",
			`FS: http.Dir("./admin/site"),`,
		}, {
			"simple /gen.go",
			"gen.go",
			`//go:generate good pages gen .`,
		}, {
			"simple app.css",
			"static/styles/app.css",
			"/* side wide CSS here */",
		}, {
			"simple site env.go",
			"site/env.go",
			"type Env struct {",
		}, {
			"simple app.js",
			"static/js/app.js",
			"function hello(name) {",
		}, {
			"simple test.txt",
			"static/public/test.txt",
			"hello world",
		}, {
			"import for page handlers.go",
			"page/handlers.go",
			`"github.com/rur/example/admin/site/site"`,
		}, {
			"content handler for page handlers.go",
			"page/handlers.go",
			`func ExampleSharedHandler(env *site.Env, rsp treetop.Response, req *http.Request) interface{} {`,
		}, {
			"base template",
			"page/templates/shared.html.tmpl",
			`<p>This template can shared by different pages</p>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, ok := gotFileMap[tt.file]
			if !ok {
				t.Errorf("Expecting file name in output: %s", tt.file)
			} else if !strings.Contains(content, tt.contains) {
				t.Errorf("Expecting %s to contain: %s, got: %s", tt.file, tt.contains, content)
			}
		})
	}
}

func filesToMap(files []File) map[string]string {
	byPath := make(map[string]string)
	for i := range files {
		file := files[i]
		byPath[filepath.Join(file.Dir, file.Name)] = string(file.Contents)
	}
	return byPath
}
