package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rur/good/generate"
)

func main() {
	switch os.Args[1] {
	case "scaffold":
		mod, err := generate.ReadGoModFile("./go.mod")
		mustNot(err)
		dest, err := generate.PrepScaffoldDir(os.Args[2])
		mustNot(err)
		files, err := generate.Scaffold(mod, dest)
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
