package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/rur/good/generate"
)

// go:embed scaffold
var scaffold embed.FS

func main() {
	switch os.Args[1] {
	case "scaffold":
		mod := generate.MustReadGoModFile("./go.mod")
		fmt.Printf("Good scaffold for %s! go version %s", mod.Module, mod.Version)
	case "page":
		mod := generate.MustReadGoModFile("./go.mod")
		fmt.Printf("Good page for %s! go version %s", mod.Module, mod.Version)
	case "routes":
		fmt.Println("Good routes")
	default:
		log.Fatalf("Unknown command %s", os.Args[1])
	}
	fmt.Println()
}

func mustBeNil(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
