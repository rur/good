package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
	"github.com/rur/good/generate"
	"github.com/rur/good/routemap"
)

var (
	usage = `usage: good <command> [<args>]

These are scaffolding commands for the Good tool:

Commands
	scaffold <site_pkg> [<page_name>...]    Create a new site scaffold at at a package relative to the working dir
	page     <site_pkg> <page_name>         Add a new page to an existing scaffold
	pages    <site_pkg>                    Generate a routes.go file from a TOML config
	routes   <page_pkg>                    Generate a routes.go file from a TOML config

`
	scaffoldUsage = `usage: good scaffold <site_pkg_rel> [<page_name>...]

Create a new site scaffold in the current Golang project

Example
	good scaffold ./admin/site dashboard settings

Arguments
	site_pkg_rel    relative import path of the new site from the current Go module
	page_name       optional list of page names to be initialized along with the site, default is 'example'

`
	pageUsage = `usage: good page <site_pkg_rel> <page_name>

Add a new page to an existing scaffold site.

Example
	good page ./admin/site settings

Arguments
	site_pkg_rel   relative import path of an existing scaffold site from the current Go module
	page_name      package name of the new page to initialize

`
	pagesUsage = `usage: good pages <site_pkg_rel>

Scan site page/ folder and update the pages.go file with any pages added manually.

Example
	good pages ./admin/site

Arguments
	site_pkg_rel   relative import path of an existing scaffold site from the current Go module

`
	routesUsage = `usage: good routes <page_pkg_rel>

Generate golang code for the routemap.toml config file in a target page. This will also populate code
for any handlers or templates that are missing from the config.

Example
	good routes ./admin/site/page/example

Arguments
	page_pkg_rel   page import path from the root of the Go module

`
)

//go:embed scaffold
var scaffold embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		log.Fatalf("Missing <command>")
	}
	switch os.Args[1] {
	case "scaffold":
		if len(os.Args) < 3 {
			fmt.Println(scaffoldUsage)
			log.Fatalf("Missing target site package path")
		}
		scaffoldCmd(os.Args[2], os.Args[3:])

	case "page":
		if len(os.Args) < 4 {
			fmt.Println(pageUsage)
			log.Fatalf("Missing required arguments")
		}
		pageCmd(os.Args[2], os.Args[3])

	case "pages":
		if len(os.Args) < 3 {
			fmt.Println(pagesUsage)
			log.Fatalf("Missing target site package path")
		}
		pagesCmd(os.Args[2])

	case "routes":
		if len(os.Args) < 3 {
			fmt.Println(routesUsage)
			log.Fatalln("Missing target scaffold page path")
		}
		routesCmd(os.Args[2])

	default:
		fmt.Println(usage)
		log.Fatalf("Unknown command %s", os.Args[1])
	}
	fmt.Println()
}

// scaffoldCmd creates a full site scaffold at a location relative to the
// current golang module.
func scaffoldCmd(sitePkgRel string, pages []string) {
	if len(pages) == 0 {
		// if no page names were listed, add a page called 'example'
		pages = []string{"example"}
	}

	// use current package to find go module
	curPkg, err := generate.GoListPackage(".")
	mustNot(err)
	sitePkg, err := generate.ParseSitePackage(curPkg.Module, os.Args[2])
	mustNot(err)
	err = generate.ValidateScaffoldLocation(sitePkg.Dir, scaffold)
	mustNot(err)
	files, err := generate.SiteScaffold(sitePkg, scaffold)
	mustNot(err)
	pageFile, err := generate.PagesScaffold(sitePkg, pages, scaffold)
	mustNot(err)
	files = append(files, pageFile)

	welcome := os.DirFS(filepath.Join("starter", "welcome"))

	for _, page := range pages {
		err = generate.ValidatePageName(page)
		mustNot(err)
		pFiles, err := generate.PageScaffold(sitePkg, page, scaffold, welcome)
		mustNot(err)
		files = append(files, pFiles...)
	}

	// FS operations
	err = generate.FlushFiles(sitePkg.Dir, files)
	mustNot(err)

	stdout, err := generate.GoFormat(sitePkg.ImportPath + "/...")
	if err != nil {
		log.Fatalf("Page '%s' scaffold was create with formatting errors: %s", sitePkg, err)
	}
	if len(stdout) > 0 {
		fmt.Println("Output from go fmt:")
		fmt.Println(stdout)
	}
	fmt.Printf("Created good scaffold for %s!", sitePkg.ImportPath)
}

// pageCmd attempts to add a new page to an existing scaffold site
func pageCmd(sitePkgRel, pageName string) {
	err := generate.ValidatePageName(pageName)
	mustNot(err)
	sitePkg, err := generate.GoListPackage(sitePkgRel)
	mustNot(err)
	err = generate.ValidatePageLocation(filepath.Join(sitePkg.Dir, "page", pageName), scaffold)
	mustNot(err)
	welcome := os.DirFS(filepath.Join("starter", "welcome"))
	files, err := generate.PageScaffold(sitePkg, pageName, scaffold, welcome)
	mustNot(err)
	err = generate.FlushFiles(sitePkg.Dir, files)
	mustNot(err)

	pageImport := fmt.Sprintf("%s/page/%s", sitePkg.ImportPath, pageName)
	fmt.Printf("Created page at %s\n", pageImport)

	pageList, err := generate.ScanSitemap(sitePkg)
	mustNot(err)
	pages, err := generate.PagesScaffold(sitePkg, pageList, scaffold)
	mustNot(err)
	err = generate.FlushFiles(sitePkg.Dir, []generate.File{pages})
	mustNot(err)

	stdout, err := generate.GoFormat(sitePkg.ImportPath + "/...")
	if err != nil {
		log.Fatalf("Page '%s' scaffold was create with fmt error: %s", pageImport, err)
	}
	if len(stdout) > 0 {
		fmt.Println("Output from go fmt:")
		fmt.Println(stdout)
	}
	fmt.Printf("Created good page for %s!", pageName)
}

// pagesCmd generates a new pages.go file by scanning the [site]/page/* directory
// for pages
func pagesCmd(sitePkgRel string) {
	sitePkg, err := generate.GoListPackage(sitePkgRel)
	mustNot(err)
	pageList, err := generate.ScanSitemap(sitePkg)
	mustNot(err)
	pages, err := generate.PagesScaffold(sitePkg, pageList, scaffold)
	mustNot(err)
	err = generate.FlushFiles(sitePkg.Dir, []generate.File{pages})
	mustNot(err)
	stdout, err := generate.GoFormat(sitePkg.ImportPath)
	if err != nil {
		log.Fatalf("Pages file for '%s' scaffold was create with formatting error: %s", sitePkg.ImportPath, err)
	}
	if len(stdout) > 0 {
		fmt.Println("Output from go fmt:")
		fmt.Println(stdout)
	}
	fmt.Printf("Updated pages.go for scaffold %s!", sitePkg.ImportPath)
}

// routesCmd will parse a routemap.toml file and generate routes, handlers and templates
// as needed
func routesCmd(pagePkgRel string) {
	pkg, err := generate.GoListPackage(pagePkgRel)
	mustNot(err)
	routemapContent, err := ioutil.ReadFile(filepath.Join(pkg.Dir, "routemap.toml"))
	mustNot(err)
	tree, err := toml.LoadBytes(routemapContent)
	mustNot(err)
	pageName := pkg.Name()
	sitePkg, err := generate.SiteFromPagePackage(pkg)
	mustNot(err)
	pageRoutes, missTpl, missHlr, err := routemap.ProcessRoutemap(
		tree,
		filepath.Join("page", pageName, "templates"),
	)
	mustNot(err)
	entries, routes, handlers, templates, err := routemap.TemplateDataForRoutes(pageRoutes, missTpl, missHlr)
	mustNot(err)
	files, err := generate.RoutesScaffold(
		sitePkg,
		pageName,
		entries,
		routes,
		handlers,
		templates,
		scaffold,
	)
	mustNot(err)
	if len(missHlr)+len(missTpl) > 0 {
		modifiedContent, err := routemap.ModifiedRoutemap(bytes.NewReader(routemapContent), missTpl, missHlr)
		mustNot(err)
		files = append(files, generate.File{
			Dir:       filepath.Join("page", pageName),
			Name:      "routemap.toml",
			Contents:  modifiedContent,
			Overwrite: true,
		})
	}

	// write files to disk
	err = generate.FlushFiles(sitePkg.Dir, files)
	mustNot(err)
	stdout, err := generate.GoFormat(pkg.ImportPath)
	if err != nil {
		log.Fatalf("Routes for '%s' page were create with formatting error: %s", pkg.ImportPath, err)
	}
	if len(stdout) > 0 {
		fmt.Println("Output from go fmt:")
		fmt.Println(stdout)
	}
	fmt.Printf("Updated routes.go for scaffold page %s!", pkg.ImportPath)
}

func mustNot(err error) {
	if err != nil {
		panic(err)
	}
}
