package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/rur/treetop"

	"github.com/rur/good/_baseline/page/page"
	"github.com/rur/good/_baseline/page/service"
)

var (
	// CLI Flags
	port    uint
	devMode bool

	//go:embed static/public
	//go:embed static/js
	//go:embed static/styles
	assets embed.FS

	//go:embed page/*/templates
	//go:embed page/templates
	templates embed.FS

	// used to serve static assets
	staticFS http.FileSystem

	// binding views to template files
	exec treetop.ViewExecutor

	env *service.Env
)

func init() {
	// CLI
	flag.BoolVar(&devMode, "dev", false, "Development mode, disable caching and other production optimizations")
	flag.UintVar(&port, "port", 8000, "Port number to bind to")
	flag.Parse()

	if devMode {
		// all static assets and templates should be read from disk at runtime
		fmt.Println("Server running in DEVELOPMENT MODE")
		staticFS = http.Dir("static")
		exec = &treetop.DeveloperExecutor{ // force templates to be re-parsed for every request
			ViewExecutor: &treetop.FileSystemExecutor{
				FS: http.Dir("./_baseline/page"), // read templates from file system
			},
		}
	} else {
		// static assets and templates are embedded at compile time
		dir, _ := fs.Sub(assets, "static")
		staticFS = http.FS(dir)
		exec = &treetop.FileSystemExecutor{
			FS: http.FS(templates),
		}
	}
}

func main() {
	// Initialize Env instance to be shared with all handlers
	env = &service.Env{
		// EDITME: initialize site-wide stuff here
		DB: nil,
	}

	m := &http.ServeMux{}

	// see ./pages.go
	registerPages(&page.DefaultHelper{
		Env: env,
		Mux: m,
	}, exec)

	if errs := exec.FlushErrors(); len(errs) > 0 {
		// templates referred by the router but could not be found
		log.Fatalf("Template errors:\n%s", errs)
	}

	m.Handle("/styles/", http.FileServer(staticFS))
	m.Handle("/js/", http.FileServer(staticFS))
	m.Handle("/public/", http.FileServer(staticFS))
	// TODO: embed treetop client in serverside library
	// m.Handle("/js/treetop.js", treetop.ServeClientLibrary)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting github.com/rur/good/_baseline/page server at %s\n", addr)

	// Bind to an addr and pass our router in
	log.Fatal(http.ListenAndServe(addr, m))
}
