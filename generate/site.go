package generate

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// SiteScaffold will return a list of files that need to be created
func SiteScaffold(mod string, dest string, pages []string, scaffold fs.FS) (files []File, err error) {
	data := struct {
		Namespace string
		Pages     []string
	}{
		Namespace: mod,
		Pages:     pages,
	}

	// main.go
	files = append(files, File{
		Dir:      dest,
		Name:     "main.go",
		Contents: mustExecute("scaffold/main.go.tmpl", data, scaffold),
	})
	// pages.go
	files = append(files, File{
		Dir:      dest,
		Name:     "pages.go",
		Contents: mustExecute("scaffold/pages.go.tmpl", data, scaffold),
	})
	// gen.go
	files = append(files, File{
		Dir:      dest,
		Name:     "gen.go",
		Contents: mustExecute("scaffold/gen.go.tmpl", data, scaffold),
	})
	// static/*
	if err = fs.WalkDir(scaffold, "scaffold/static", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, File{
			Dir:      filepath.Join(dest, strings.TrimPrefix(filepath.Dir(path), "scaffold/")),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	}); err != nil {
		return
	}
	// app/*
	if err = fs.WalkDir(scaffold, "scaffold/app", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, File{
			Dir:      filepath.Join(dest, strings.TrimPrefix(filepath.Dir(path), "scaffold/")),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	}); err != nil {
		return
	}
	// page/
	err = fs.WalkDir(scaffold, "scaffold/page", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			if path == "scaffold/page/name" {
				return fs.SkipDir
			}
			return nil
		}
		files = append(files, File{
			Dir:      filepath.Join(dest, strings.TrimPrefix(filepath.Dir(path), "scaffold/")),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	})
	return
}

// ValidateScaffoldPackage will check that a relative go package path is available for
// a site scaffold to to be installed and return the go import path and directory for the
// new scaffold site.
//
// Note that '.' will attempt to install the site in the root directory of the current go module
//
// If the target directory is not empty, this will scan for conflicts against the scaffold
func ValidateScaffoldPackage(pkg GoModule, name string, scaffold fs.FS) (string, string, error) {
	var (
		sitePkg, siteDir string
	)
	if name == "." {
		sitePkg = pkg.Path
		siteDir = ""
	} else if strings.HasPrefix(name, pkg.Path) {
		return "", "", fmt.Errorf("site package name must be relative to the current module, got %s", name)
	} else {
		// strip relative prefix since being relative is assumed
		name = strings.TrimPrefix(name, "./")
		parts := strings.Split(name, "/")
		sitePkg = strings.Join([]string{pkg.Path, name}, "/")
		siteDir = filepath.Join(parts...)
	}

	// now try to check if there will be files write conflicts
	// build block list index
	blocked := struct{}{} // zero size sentinel
	blockList := make(map[string]struct{})
	entries, err := fs.ReadDir(scaffold, "scaffold")
	if err != nil {
		return "", "", fmt.Errorf("failed to open scaffold folder: %s", err)
	}
	for _, entry := range entries {
		blockList[strings.TrimSuffix(entry.Name(), ".tmpl")] = blocked
	}

	// Scan for conflict between the scaffold and the target FS
	// As a sanity check, accept at most 500 dept one child names
	fh, err := os.Open(siteDir)
	if err == os.ErrNotExist {
		// this is fine
		return sitePkg, siteDir, nil
	}
	defer fh.Close()
	names, err := fh.Readdirnames(500)
	if err != io.EOF {
		for i := 0; i < len(names); i++ {
			if _, ok := blockList[names[i]]; ok {
				return "", "", fmt.Errorf("conflicting file or direcotry '%s'", names[i])
			}
		}
	}
	return sitePkg, siteDir, nil
}

// mustExecute will execute a template against data or panic!
// Since the templates are embedded we can treat failure at this stage
// as a bug
func mustExecute(name string, data interface{}, scaffold fs.FS) []byte {
	tmpl, err := template.New(path.Base(name)).Delims("<<", ">>").ParseFS(scaffold, name)
	if err != nil {
		log.Fatalln("Failed to parse template", name, err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Fatalln("Failed to execute template", name, err)
	}
	return buf.Bytes()
}
