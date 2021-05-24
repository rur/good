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
		pkg, err := generate.GoListPackage(".")
		mustNot(err)
		siteDir, siteModule, err := generate.ValidateScaffoldPackage(pkg.Module, os.Args[2], scaffold)
		mustNot(err)
		files, err := generate.SiteScaffold(siteModule, siteDir, scaffold)
		mustNot(err)

		// setup initial pages
		var pages []string
		if len(os.Args) > 2 {
			pages = os.Args[2:]
		} else {
			// if no page names were listed, add a page called 'example'
			pages = []string{"example"}
		}
		for _, page := range pages {
			pFiles, err := generate.ScaffoldPage(pkg.Module.Path, os.Args[2], page, scaffold)
			mustNot(err)
			files = append(files, pFiles...)
		}

		// TODO: run pages update on scaffold

		err = generate.FlushFiles(files)
		mustNot(err)
		sitePkg, err := generate.GoListPackage("./" + os.Args[2])
		if err != nil {
			log.Fatalf("Scaffold was create with errors: %s", err)
		}
		fmt.Printf("Created good scaffold for %s!", sitePkg.ImportPath)

	case "page":
		pkg, err := generate.GoListPackage(".")
		mustNot(err)
		mustNot(err)
		fmt.Printf("Good page for %s!", pkg.ImportPath)
	case "routes":
		fmt.Println("Good routes")
	default:
		log.Fatalf("Unknown command %s", os.Args[1])
	}
	fmt.Println()
}

func mustNot(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}