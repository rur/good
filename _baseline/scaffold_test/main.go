package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rur/treetop"

	"github.com/rur/good/baseline/scaffold_test/page"
	"github.com/rur/good/baseline/scaffold_test/site"
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
)

func init() {
	flag.UintVar(&port, "port", 8000, "Port number to bind to")
	flag.BoolVar(&devMode, "dev", false, "Development mode, disable caching and other production optimizations")
	flag.Parse()
}

func main() {
	var (
		staticFS http.FileSystem
		exec     treetop.ViewExecutor
	)
	if devMode {
		// all static assets and templates should be read from disk at runtime
		fmt.Println("Server running in DEVELOPMENT MODE")
		staticFS = http.Dir(filepath.Join("./baseline/scaffold_test", "static"))
		exec = &treetop.DeveloperExecutor{ // force templates to be re-parsed for every request
			ViewExecutor: &treetop.FileSystemExecutor{
				// this assumes you are runing the dev server from your project root
				FS:          http.Dir("./baseline/scaffold_test"),
				KeyedString: page.KeyedTemplates,
			},
		}
	} else {
		// static assets and templates are embedded at compile time
		dir, _ := fs.Sub(assets, "static")
		staticFS = http.FS(dir)
		exec = &treetop.FileSystemExecutor{
			FS:          http.FS(templates),
			KeyedString: page.KeyedTemplates,
		}
	}

	// Initialize Env singleton, it will be bound to handlers using the page.BindEnv helper
	env := &site.Env{
		// EDITME: initialize site-wide stuff here, for example...
		HTTPS:    !devMode,
		ErrorLog: log.New(os.Stderr, "[error]: ", log.Llongfile),
		WarnLog:  log.New(os.Stdout, "[warn]: ", log.Llongfile),
		InfoLog:  log.New(os.Stdout, "[info]: ", log.Llongfile),
		DB:       nil,
	}

	// EDITME: it is recommended to replace this with your preferred routing library
	m := &http.ServeMux{}

	// see ./pages.go
	registerPages(page.Helper{
		Env: env,
		Mux: m,
	}, exec)

	if errs := exec.FlushErrors(); len(errs) > 0 {
		// flush any template related errors from the routes config
		// templates are eagerly loaded and parsed to try surface issues at startup
		log.Fatalf("Template errors:\n%s", errs)
	}

	{
		// static files
		public, _ := fs.Sub(assets, "static/public")
		m.Handle("/favicon.ico", http.FileServer(http.FS(public)))
		m.Handle("/humans.txt", http.FileServer(http.FS(public)))
		m.Handle("/js/treetop.js", treetop.ServeClientLibrary)

		// static folders
		m.Handle("/js/", http.FileServer(staticFS))
		m.Handle("/styles/", http.FileServer(staticFS))
		m.Handle("/public/", http.FileServer(staticFS))
	}

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting github.com/rur/good/baseline/scaffold_test server at %s\n", addr)

	// Bind to an addr and pass our router in
	log.Fatal(http.ListenAndServe(addr, m))
}
