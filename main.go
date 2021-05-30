package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/rur/good/generate"
)

var (
	usage = `usage: good <command> [<args>]

These are scaffolding commands for the Good tool:

	scaffold <package_name> [<page_name>...]    Create a new site scaffold at at a package relative to the working dir
	page     <package_name> <page_name>         Add a new page to an existing scaffold
	routes   <routemap_toml>                    Generate a routes.go file from a TOML config

`
	scaffoldUsage = `usage: good scaffold <site_pkg_rel> [<page_name>...]

Create a new site scaffold in the current Golang project

eg.
	good scaffold ./admin/site dashboard settings

site_pkg_rel    import path of the new site, relative to the root of the current Go module
page_name       optional list of page names to be initialized along with the site, default is 'example'

`
	pageUsage = `usage: good page <site_pkg_rel> <page_name>

Add a new page to an existing scaffold site.

eg.
	good page ./admin/site settings

site_pkg_rel   import path of an existing _good scaffold_ site
page_name      package name of the new page to initialize

`
	routesUsage = `usage: good routes <page_pkg_rel>

Generate golang code for the routing config in a target page and populate code for any handlers or templates
that are missing.

eg.
	good routes ./admin/site/page/example

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
		var pages []string
		if len(os.Args) < 3 {
			fmt.Println(scaffoldUsage)
			log.Fatalf("Missing target site package path")
		} else if len(os.Args) > 3 {
			pages = os.Args[3:]
		} else {
			// if no page names were listed, add a page called 'example'
			pages = []string{"example"}
		}

		pkg, err := generate.GoListPackage(".")
		mustNot(err)
		sitePkg, siteDir, err := generate.ValidateScaffoldPackage(pkg.Module, os.Args[2], scaffold)
		mustNot(err)
		files, err := generate.SiteScaffold(sitePkg, siteDir, pages, scaffold)
		mustNot(err)

		for _, page := range pages {
			pFiles, err := generate.ScaffoldPage(sitePkg, siteDir, page, scaffold)
			mustNot(err)
			files = append(files, pFiles...)
		}
		// FS operations
		err = generate.FlushFiles(pkg.Module.Dir, files)
		mustNot(err)

		if err := generate.GoFormat(sitePkg + "/..."); err != nil {
			log.Fatalf("Scaffold was create with errors: %s", err)
		}
		fmt.Printf("Created good scaffold for %s!", sitePkg)

	case "page":
		if len(os.Args) < 4 {
			fmt.Println(pageUsage)
			log.Fatalf("Missing required arguments")
		}
		pkg, err := generate.GoListPackage(".")
		mustNot(err)
		sitePkg, siteDir, err := generate.ParseSitePackage(pkg.Module, os.Args[2])
		mustNot(err)
		files, err := generate.ScaffoldPage(sitePkg, siteDir, os.Args[3], scaffold)
		mustNot(err)
		// FS operations
		err = generate.FlushFiles(pkg.Module.Dir, files)
		mustNot(err)

		pagePkg := fmt.Sprintf("%s/page/%s", sitePkg, os.Args[3])
		if err := generate.GoFormat(pagePkg + "/..."); err != nil {
			log.Fatalf("Page '%s' scaffold was create with errors: %s", pagePkg, err)
		}
		fmt.Printf("Created good page for %s!", pagePkg)
	case "routes":
		fmt.Println(routesUsage)
		log.Fatalf("Good routes is not implemented yet!")

	default:
		fmt.Println(usage)
		log.Fatalf("Unknown command %s", os.Args[1])
	}
	fmt.Println()
}

func mustNot(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
