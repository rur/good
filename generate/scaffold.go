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
func SiteScaffold(pkg GoPackage, scaffold fs.FS) (files []File, err error) {
	dest, err := pkg.RelPath()
	if err != nil {
		return
	}
	data := struct {
		SiteDirRel string
		Namespace  string
	}{
		SiteDirRel: dest,
		Namespace:  pkg.ImportPath,
	}

	// main.go
	files = append(files, File{
		Dir:      "",
		Name:     "main.go",
		Contents: mustExecute("scaffold/main.go.tmpl", data, scaffold),
	})
	// gen.go
	files = append(files, File{
		Dir:      "",
		Name:     "gen.go",
		Contents: mustExecute("scaffold/gen.go.tmpl", data, scaffold),
	})
	// README.md
	files = append(files, File{
		Dir:      "",
		Name:     "README.md",
		Contents: mustExecute("scaffold/README.md.tmpl", data, scaffold),
	})
	// docs/*
	if err = fs.WalkDir(scaffold, "scaffold/docs", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		files = append(files, File{
			Dir:      strings.TrimPrefix(filepath.Dir(path), "scaffold"+string(os.PathSeparator)),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	}); err != nil {
		return
	}
	// static/*
	if err = fs.WalkDir(scaffold, "scaffold/static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		files = append(files, File{
			Dir:      strings.TrimPrefix(filepath.Dir(path), "scaffold"+string(os.PathSeparator)),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	}); err != nil {
		return
	}
	// service/*
	if err = fs.WalkDir(scaffold, "scaffold/service", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		files = append(files, File{
			Dir:      strings.TrimPrefix(filepath.Dir(path), "scaffold"+string(os.PathSeparator)),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	}); err != nil {
		return
	}
	// page/
	err = fs.WalkDir(scaffold, "scaffold/page", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path == "scaffold/page/default" || path == "scaffold/page/name" {
				return fs.SkipDir
			}
			return nil
		}
		files = append(files, File{
			Dir:      strings.TrimPrefix(filepath.Dir(path), "scaffold"+string(os.PathSeparator)),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	})
	return
}

// ParseSitePackage will normalize a relative path from a valid Go module and
// return a golang package import path and directory path relative to the module dir
func ParseSitePackage(mod GoModule, name string) (sitePkg GoPackage, err error) {
	var (
		sitePath, siteDir string
	)
	if name == "." {
		sitePath = mod.Path
		siteDir = ""
	} else if strings.HasPrefix(name, mod.Path) {
		err = fmt.Errorf("site package name must be relative to the current module, got %s", name)
		return
	} else {
		// strip relative prefix since being relative is assumed
		name = strings.TrimPrefix(name, "./")
		parts := strings.Split(name, "/")
		sitePath = strings.Join([]string{mod.Path, name}, "/")
		siteDir = filepath.Join(parts...)
	}
	sitePkg = GoPackage{
		Dir:        filepath.Join(mod.Dir, siteDir),
		ImportPath: sitePath,
		Module:     mod,
	}
	return
}

// ValidateScaffoldPackage will check that a relative go package path is available for
// a site scaffold to to be installed and return the go import path and directory for the
// new scaffold site.
//
// Note that '.' will attempt to install the site in the root directory of the current go module
//
// If the target directory is not empty, this will scan for conflicts against the scaffold
func ValidateScaffoldLocation(siteDir string, scaffold fs.FS) error {
	// now try to check if there will be files write conflicts
	// build block list index
	blocked := struct{}{} // zero size sentinel
	blockList := make(map[string]struct{})
	entries, err := fs.ReadDir(scaffold, "scaffold")
	if err != nil {
		return fmt.Errorf("failed to open scaffold folder: %s", err)
	}
	for _, entry := range entries {
		blockList[strings.TrimSuffix(entry.Name(), ".tmpl")] = blocked
	}

	// Scan for conflict between the scaffold and the target FS
	// As a sanity check, accept at most 500 dept one child names
	fh, err := os.Open(siteDir)
	if os.IsNotExist(err) {
		// no destination folder, no conflicts, all good
		return nil
	} else if err != nil {
		return fmt.Errorf("error while scanning target dir: %s", err)
	}
	defer fh.Close()
	names, err := fh.Readdirnames(500)
	if err == io.EOF {
		// empty dir, no conflicts, all good
		return nil
	} else {
		for i := 0; i < len(names); i++ {
			if _, ok := blockList[names[i]]; ok {
				return fmt.Errorf("conflicting file or direcotry '%s'", names[i])
			}
		}
	}
	return nil
}

// mustExecute will execute a template against data or panic!
// Since the templates are embedded we can treat failure at this stage
// as a bug
func mustExecute(name string, data interface{}, scaffold fs.FS) []byte {
	bites, err := tryExecute(name, data, scaffold)
	if err != nil {
		log.Fatalln("Failed to parse template", name, err)
	}
	return bites
}

// tryExecute will attempt to load a template file form and FS and then
// execute the template using the provided data.
// Any error that occurs will be returned
func tryExecute(name string, data interface{}, scaffold fs.FS) ([]byte, error) {
	tmpl, err := template.New(path.Base(name)).Delims("[#", "#]").ParseFS(scaffold, name)
	if err != nil {
		log.Fatalln("Failed to parse template", name, err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
