package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/rur/good/generate"
)

var usage = `usage: good <command> [<args>]

These are scaffolding commands for the Good tool:

	scaffold <package_name> [<page_name>...]    Create a new site scaffold at at a package relative to the working dir
	page     <package_name> <page_name>         Add a new page to an existing scaffold
	routes   <routemap_toml>                    Generate a routes.go file from a TOML config

TODO: add more docs

`

//go:embed scaffold
var scaffold embed.FS

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		log.Fatalf("Missing command")
	}
	switch os.Args[1] {
	case "scaffold":
		// setup initial pages
		var pages []string
		if len(os.Args) > 3 {
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
		if len(os.Args) < 3 {
			fmt.Println(usage)
			log.Fatalf("Missing page name")
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
		fmt.Println("Good routes")
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
