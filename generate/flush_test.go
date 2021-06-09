package generate

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestFlushFiles(t *testing.T) {
	dir := os.TempDir()

	tests := []struct {
		name      string
		files     []File
		wantFiles map[string][]string
		wantErr   bool
	}{
		{
			name: "basic",
			files: []File{
				{
					Dir:      "",
					Name:     "test.txt",
					Contents: []byte(`This is a test`),
				},
			},
			wantFiles: map[string][]string{
				filepath.Join("test.txt"): {
					"This is a test",
				},
			},
		},
		{
			name: "files_with_dirs",
			files: []File{
				{
					Dir:      "",
					Name:     "test.txt",
					Contents: []byte(`This is a test`),
				},
				{
					Dir:      filepath.Join("a", "b", "c"),
					Name:     "test.txt",
					Contents: []byte(`This is a test ABC`),
				},
				{
					Dir:      filepath.Join("a", "b"),
					Name:     "test.txt",
					Contents: []byte(`This is a test AB`),
				},
				{
					Dir:      filepath.Join("a", "b"),
					Name:     "test_2.txt",
					Contents: []byte(`This is a test 2 AB`),
				},
			},
			wantFiles: map[string][]string{
				filepath.Join("test.txt"):                {"This is a test"},
				filepath.Join("a", "b", "c", "test.txt"): {"This is a test", "ABC"},
				filepath.Join("a", "b", "test.txt"):      {"This is a test", "AB"},
				filepath.Join("a", "b", "test_2.txt"):    {"This is a test", "2 AB"},
			},
		},
		{
			name: "overwrite",
			files: []File{
				{
					Dir:  "",
					Name: "test.txt",
					Contents: []byte(`This is a test
					asdf
					asdf
					asdf
					asd
					fsa
					d
					`),
				},
				{
					Dir:      "",
					Name:     "test.txt",
					Contents: []byte(`This is a new version`),
				},
			},
			wantFiles: map[string][]string{
				filepath.Join("test.txt"): {
					"^This is a new version$",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testdir := filepath.Join(dir, tt.name)
			err := FlushFiles(testdir, tt.files)
			if (err != nil) != tt.wantErr {
				t.Errorf("FlushFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			generatedFS := os.DirFS(testdir)
			gotFileMap := make(map[string]string)
			err = fs.WalkDir(generatedFS, ".", func(path string, d fs.DirEntry, err error) error {
				if d != nil && !d.IsDir() {
					content, err := fs.ReadFile(generatedFS, path)
					if err != nil {
						return err
					}
					gotFileMap[path] = string(content)
				}
				return nil
			})
			if err != nil {
				t.Error(err)
				return
			}

			for file, checks := range tt.wantFiles {
				content, ok := gotFileMap[file]
				if !ok {
					t.Errorf("Expecting file '%s'", file)
				} else {
					for _, pattern := range checks {
						if !regexp.MustCompile(pattern).MatchString(content) {
							t.Errorf("Expecting file '%s' to match %s\nGOT: %s", file, pattern, content)
						}
					}
				}
			}
			for file := range gotFileMap {
				_, ok := tt.wantFiles[file]
				if !ok {
					t.Errorf("Unexpected file '%s'", file)
				}
			}

		})
	}
}
