package generate

import (
	"embed"
	"io/fs"
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
			wantErr:  "conflicting file or direcotry 'file.go'",
		},
		{
			name:     "conflicting directory",
			location: "./generate/testdata/with_conflict_folder",
			wantErr:  "conflicting file or direcotry 'folder'",
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
		name    string
		input   string
		wantPkg string
		wantDir string
		wantErr string
	}{
		{
			name:    "basic",
			input:   "./admin/site",
			wantPkg: "github.com/rur/good/admin/site",
			wantDir: filepath.Join("admin", "site"),
		},
		{
			name:    "using . as destination",
			input:   ".",
			wantPkg: "github.com/rur/good",
			wantDir: "",
		},
		{
			name:    "embedded import",
			input:   "github.com/rur/good/admin/site",
			wantErr: "site package name must be relative to the current module, got github.com/rur/good/admin/site",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSitePkg, gotSiteDir, err := ParseSitePackage(pkg.Module, tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("ParseSitePackage() expecting an error containing message %s", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("ParseSitePackage() expecting error to contain '%s', got '%s'", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("ParseSitePackage() unexpected error = %s", err)
			} else {
				if gotSitePkg != tt.wantPkg {
					t.Errorf("ParseSitePackage() gotSitePkg = %v, want %v", gotSitePkg, tt.wantPkg)
				}
				if gotSiteDir != tt.wantDir {
					t.Errorf("ParseSitePackage() gotSiteDir = %v, want %v", gotSiteDir, tt.wantDir)
				}
			}
		})
	}
}
