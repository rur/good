package generate

import (
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestGoListPackageError(t *testing.T) {
	// fs error
	if _, err := GoListPackage("./fake"); err == nil {
		t.Error("Expecting error for fake directory")
	} else if !strings.Contains(err.Error(), "failed to load a go package for path './fake'") {
		t.Errorf("Expecting reason for failure, got %s", err)
	}
}

func TestGoListPackage(t *testing.T) {
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
	tests := []struct {
		name           string
		input          string
		wantDir        string
		wantImport     string
		wantModulePath string
		wantModuleDir  string
		wantName       string
		wantErr        bool
	}{
		{
			name:           "auto dot",
			input:          "",
			wantDir:        "^.+/good/generate$",
			wantImport:     "^github.com/rur/good/generate$",
			wantModulePath: "^github.com/rur/good$",
			wantName:       "generate",
			wantModuleDir:  "^.+/good$",
		},
		{
			name:           "current package",
			input:          ".",
			wantDir:        "^.+/good/generate$",
			wantImport:     "^github.com/rur/good/generate$",
			wantModulePath: "^github.com/rur/good$",
			wantName:       "generate",
			wantModuleDir:  "^.+/good$",
		},
		{
			name:           "current package",
			input:          "../",
			wantDir:        "^.+/good$",
			wantImport:     "^github.com/rur/good$",
			wantModulePath: "^github.com/rur/good$",
			wantName:       "main",
			wantModuleDir:  "^.+/good$",
		},
		{
			name:    "missing dir",
			input:   "./fake",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg, err := GoListPackage(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoListPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				if !assertPath("Dir", pkg.Dir, tt.wantDir) {
					return
				}
				if !assertPattern("ImportPath", pkg.ImportPath, tt.wantImport) {
					return
				}
				if !assertPattern("Module.Path", pkg.Module.Path, tt.wantModulePath) {
					return
				}
				if !assertPath("Module.Dir", pkg.Module.Dir, tt.wantModuleDir) {
					return
				}
				if pkg.Name != tt.wantName {
					t.Errorf("expecting name %q, got %q", tt.wantName, pkg.Name)
				}
			}
		})
	}
}
