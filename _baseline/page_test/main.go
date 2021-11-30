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

	"github.com/rur/good/baseline/page_test/page"
	"github.com/rur/good/baseline/page_test/site"
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

	// site environment singleton
	env *site.Env
)

func init() {
	// CLI
	flag.BoolVar(&devMode, "dev", false, "Development mode, disable caching and other production optimizations")
	flag.UintVar(&port, "port", 8000, "Port number to bind to")
	flag.Parse()

	if devMode {
		// all static assets and templates should be read from disk at runtime
		fmt.Println("Server running in DEVELOPMENT MODE")
		staticFS = http.Dir(filepath.Join("./baseline/page_test", "static"))
		exec = &treetop.DeveloperExecutor{ // force templates to be re-parsed for every request
			ViewExecutor: &treetop.FileSystemExecutor{
				// this assumes you are runing the dev server from your project root
				FS:          http.Dir("./baseline/page_test"),
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

	// Initialize Env instance to be shared with all handlers
	env = &site.Env{
		// EDITME: initialize site-wide stuff here, for example...
		HTTPS:    !devMode,
		ErrorLog: log.New(os.Stderr, "[error]: ", log.Llongfile),
		WarnLog:  log.New(os.Stdout, "[warn]: ", log.Llongfile),
		InfoLog:  log.New(os.Stdout, "[info]: ", log.Llongfile),
		DB:       nil,
	}
}

func main() {
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
	fmt.Printf("Starting github.com/rur/good/baseline/page_test server at %s\n", addr)

	// Bind to an addr and pass our router in
	log.Fatal(http.ListenAndServe(addr, m))
}
