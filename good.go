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

	scaffold <package_name>
	page     <package_name> <page_name>
	routes   <page_package_name>

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
		mod, err := generate.ReadGoModFile("./go.mod")
		mustNot(err)
		dest, err := generate.ValidateScaffoldPath(os.Args[2])
		mustNot(err)
		files, err := generate.Scaffold(mod, dest, scaffold)
		mustNot(err)
		err = generate.FlushFiles(files)
		mustNot(err)
		// TODO: create pages for os.Args[3:] default to a single example page
		fmt.Printf("Created good scaffold for %s! go version %d.%d", mod.Module, mod.MajorVersion, mod.MinorVersion)

	case "page":
		mod, err := generate.ReadGoModFile("./go.mod")
		mustNot(err)
		fmt.Printf("Good page for %s! go version %d.%d", mod.Module, mod.MajorVersion, mod.MinorVersion)
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
