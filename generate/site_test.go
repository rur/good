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
	scaffold, _ := fs.Sub(testScaffold, "testdata")
	pkg, err := GoListPackage(".")
	if err != nil {
		t.Fatal("ValidateScaffoldPackage error getting local package details", err)
	}
	tests := []struct {
		name    string
		path    string
		wantPkg string
		wantDir string
		wantErr string
	}{
		{
			name:    "basic",
			path:    "./admin/site",
			wantPkg: "github.com/rur/good/admin/site",
			wantDir: filepath.Join(pkg.Module.Dir, "admin", "site"),
		},
		{
			name:    "using . as destination",
			path:    ".",
			wantPkg: "github.com/rur/good",
			wantDir: pkg.Module.Dir,
		},
		{
			name:    "conflicting file",
			path:    "./generate/testdata/with_conflict_file",
			wantErr: "conflicting file or direcotry 'file.go'",
		},
		{
			name:    "conflicting directory",
			path:    "./generate/testdata/with_conflict_folder",
			wantErr: "conflicting file or direcotry 'folder'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := ValidateScaffoldPackage(pkg.Module, tt.path, scaffold)
			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("ValidateScaffoldPackage() expecting an error containing message %s", tt.wantErr)
					return
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("ValidateScaffoldPackage() expecting error to contain '%s', got '%s'", tt.wantErr, err)
				}
			} else if err != nil {
				t.Errorf("ValidateScaffoldPackage() unexpected error = %v", err)
				return
			}
			if got != tt.wantPkg {
				t.Errorf("ValidateScaffoldPackage() got = %v, want %v", got, tt.wantPkg)
			}
			if got1 != tt.wantDir {
				t.Errorf("ValidateScaffoldPackage() got1 = %v, want %v", got1, tt.wantDir)
			}
		})
	}
}
