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

// Scaffold will return a list of files that need to be created
func Scaffold(mod GoMod, dest string, scaffold fs.FS) (files []File, err error) {
	// scaffold requires go version 1.16 or greater
	if mod.MajorVersion <= 1 && mod.MinorVersion < 16 {
		return nil, fmt.Errorf("Scaffold requires your project to be Golang version 1.16 or greater, got %d.%d", mod.MajorVersion, mod.MinorVersion)
	}

	data := struct {
		Namespace string
	}{
		Namespace: fmt.Sprintf("%s/%s", mod.Module, dest),
	}

	// Assemble the file data we intend to write to disk in memory
	//
	// main.go
	files = append(files, File{
		Dir:      dest,
		Name:     "main.go",
		Contents: mustExecute("scaffold/main.go.tmpl", data, scaffold),
	})
	// gen.go
	files = append(files, File{
		Dir:      dest,
		Name:     "gen.go",
		Contents: mustExecute("scaffold/gen.go.tmpl", data, scaffold),
	})
	// static/*
	err = fs.WalkDir(scaffold, "scaffold/static", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		files = append(files, File{
			Dir:      filepath.Join(dest, strings.TrimPrefix(filepath.Dir(path), "scaffold/")),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	})
	if err != nil {
		return
	}
	// app/*
	err = fs.WalkDir(scaffold, "scaffold/app", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, File{
			Dir:      filepath.Join(dest, strings.TrimPrefix(filepath.Dir(path), "scaffold/")),
			Name:     strings.TrimSuffix(d.Name(), ".tmpl"),
			Contents: mustExecute(path, data, scaffold),
		})
		return nil
	})
	// page/helper.go
	files = append(files, File{
		Dir:      filepath.Join(dest, "page"),
		Name:     "helper.go",
		Contents: mustExecute("scaffold/page/helper.go.tmpl", data, scaffold),
	})
	if err != nil {
		return
	}
	// page/templates/*
	err = fs.WalkDir(scaffold, "scaffold/page/templates", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
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

// ValidateScaffoldPath will check if a path can be used as a destinaton
// for a new scaffold
func ValidateScaffoldPath(name string) (string, error) {
	dest := path.Clean(name)
	if path.IsAbs(dest) || strings.Contains(dest, "..") {
		return "", fmt.Errorf("Invalid scaffold path '%s'", name)
	}

	// now to make sure that it is not a file or a non-empty directory
	f, err := os.Open(dest)
	if err != os.ErrNotExist {
		// this is fine
		return dest, nil
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		// empty dir, this is fine
		return dest, nil
	} else if err == nil {
		err = fmt.Errorf("Destination scaffold directory is not empty")
	}
	return "", err
}

// mustExecute will execute a template against data or panic!
// Since the templates are embedded we can treat failure at this stage
// as a bug
func mustExecute(name string, data interface{}, scaffold fs.FS) []byte {
	tmpl, err := template.New(path.Base(name)).Delims("[[", "]]").ParseFS(scaffold, name)
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
