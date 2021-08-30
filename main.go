package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/rur/good/generate"
	"github.com/rur/good/routemap"
)

var (
	usage = `usage: good <command> [<args>]

These are scaffolding commands for the Good tool:

Commands
	scaffold     <site_pkg>              Create a new site scaffold at at a package relative to the working dir
	page         <site_pkg> <page_name>  Add a new page to an existing scaffold
	pages gen    <site_pkg>              Re-generate the pages.go file
	pages list   <site_pkg>              List pages in a site to stdout
	pages delete <page_pkg>              Remove a page from a site
	routes       <page_pkg>              Re-generate the routes.go file for a specified scaffold page package
	starter      <out_dir>               Create a new dir and populate it with a template for a custom starter page

`
	scaffoldUsage = `usage: good scaffold <site_pkg_rel>

Create a new site scaffold in the current Golang project

Example
	good scaffold ./admin/site

Arguments
	site_pkg_rel    relative import path of the new site from the current Go module

Options
    -h
	Print usage for scaffold command

`
	pageUsage = `usage: good page <site_pkg_rel> <page_name> [--starter <path>]

Add a new page to an existing scaffold site.

Example
	good page ./admin/site settings

Arguments
	site_pkg_rel   relative import path of an existing scaffold site from the current Go module
	page_name      package name of the new page to initialize

Options
    -h
	Print usage for page command
    --starter <dir path>
        Use specified directory as the starter template for the new page scaffold.

	There are a number of built-in page starters you can choose from. These are specified using
	a single colon prefix. Here is the list of available build-ins you can try:

		:basic                 No layout, just a simple page scaffold (default)
		:bootstrap5/layout     Useful Bootsrap v5.0 web console layout
		:bootstrap5/examples   Working demos with the Bootsrap v5.0 layout and components
		:bootstrap5/login      TODO!
		:bootstrap5/datatable  TODO!
		:minimum               Most bare bones option
		:intro                 Introduction page to the good scaffold

`
	pagesUsage = `usage: good pages <site_pkg_rel>

Scan site page/ folder and update the pages.go file with any pages added manually.

Example
	good pages [-h] [command] <args>

Commands
	gen <site path>
		Generated a new pages.go file, useful when the site page folders has been modifed
	list <site path>
		List named pages found for a site to stdout
	delete <page path>
		Delete folder for a named page and regenerate the pages.go file.
		Convenience command, simalar to 'rm -rf <page path> && good pages gen <site path>'

Options
    -h
	Print usage for pages command

`
	routesUsage = `usage: good routes <page_pkg_rel>

Generate golang code for the routemap.toml config file in a target page. This will also populate code
for any handlers or templates that are missing from the config.

Example
	good routes ./admin/site/page/example

Arguments
	page_pkg_rel   page import path from the root of the Go module

Options
    -h
	Print usage for routes command

`
	starterUsage = `usage: good starter <out_dir>

Generate a template directory for a custom starter page that can be used with the
'good page' command.


Example
	$ good starter ./admin/utils/customstarter
    [...customize the files...]

    # create a new page using your starter
    $ good page ./admin/site awesomenewpage --starter ./admin/utils/customstarter

Arguments
	out_dir   a not-already-existing path where a folder will be created

Options
    -h
	Print usage for starter command

`
)

//go:embed scaffold
var scaffold embed.FS

//go:embed starter
var starter embed.FS

func main() {
	pArgs := os.Args[1:]
	fArgs := make(map[string]string)
	for i := range pArgs {
		if pArgs[i][0] == '-' {
			// treat all subsequent args as CLI flags
			for j := i; j < len(pArgs); j++ {
				// cheap and cheerful arg parsing
				if j < len(pArgs)-1 && pArgs[j+1][0] != '-' {
					fArgs[pArgs[j]] = pArgs[j+1]
					j++
				} else {
					fArgs[pArgs[j]] = ""
				}
			}
			pArgs = pArgs[:i]
			break
		}
	}

	if len(pArgs) < 1 {
		if _, help := fArgs["-h"]; help {
			fmt.Println(usage)
			return
		}
		fmt.Println(usage)
		log.Fatalf("Missing <command>")
	}
	switch pArgs[0] {
	case "scaffold":
		if _, help := fArgs["-h"]; help {
			fmt.Println(scaffoldUsage)
			return
		}
		if len(pArgs) < 2 {
			fmt.Println(scaffoldUsage)
			log.Fatalf("Missing target site package path")
		}
		scaffoldCmd(pArgs[1])

	case "page":
		if _, help := fArgs["-h"]; help {
			fmt.Println(pageUsage)
			return
		}
		if len(pArgs) < 3 {
			fmt.Println(pageUsage)
			log.Fatalf("Missing required arguments")
		}
		starterTemplatePath := fArgs["--starter"]
		fmt.Printf("Starter template %#v \n", starterTemplatePath)
		pageCmd(pArgs[1], pArgs[2], starterTemplatePath)

	case "pages":
		if _, help := fArgs["-h"]; help {
			fmt.Println(pagesUsage)
			return
		}
		if len(pArgs) < 3 {
			fmt.Println(pagesUsage)
			log.Fatalf("Incomplete pages command: %v", strings.Join(pArgs, " "))
		}
		switch pArgs[1] {
		case "gen":
			pagesCmd(pArgs[2])

		case "list":
			listPagesCmd(pArgs[2])

		case "delete":
			deletePageCmd(pArgs[2])

		default:
			fmt.Println(pagesUsage)
			log.Fatalf("Unknown pages command '%s'", pArgs[1])
		}

	case "routes":
		if _, help := fArgs["-h"]; help {
			fmt.Println(routesUsage)
			return
		}
		if len(pArgs) < 2 {
			fmt.Println(routesUsage)
			log.Fatalln("Missing target scaffold page path")
		}
		routesCmd(pArgs[1])

	case "starter":
		if _, help := fArgs["-h"]; help {
			fmt.Println(starterUsage)
			return
		}
		if len(pArgs) < 2 {
			fmt.Println(starterUsage)
			log.Fatalln("Missing target starter folder path")
		}
		starterCmd(pArgs[1])

	default:
		fmt.Println(usage)
		log.Fatalf("Unknown command %s", pArgs[0])
	}
	fmt.Println()
}

// scaffoldCmd creates a full site scaffold at a location relative to the
// current golang module.
func scaffoldCmd(sitePkgRel string) {
	// use current package to find go module
	curPkg, err := generate.GoListPackage("./...")
	mustNot(err)
	sitePkg, err := generate.ParseSitePackage(curPkg.Module, sitePkgRel)
	mustNot(err)
	err = generate.ValidateScaffoldLocation(sitePkg.Dir, scaffold)
	mustNot(err)
	files, err := generate.SiteScaffold(sitePkg, scaffold)
	mustNot(err)
	pageFile, err := generate.PagesScaffold(sitePkg, []string{"intro"}, scaffold)
	mustNot(err)
	files = append(files, pageFile)
	start, err := fs.Sub(starter, "starter/intro")
	mustNot(err)
	pFiles, err := generate.PageScaffold(sitePkg, "intro", scaffold, start)
	mustNot(err)
	files = append(files, pFiles...)

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
func pageCmd(sitePkgRel, pageName, starterTemplatePath string) {
	err := generate.ValidatePageName(pageName)
	mustNot(err)
	sitePkg, err := generate.GoListPackage(sitePkgRel)
	mustNot(err)
	err = generate.ValidatePageLocation(filepath.Join(sitePkg.Dir, "page", pageName), scaffold)
	mustNot(err)
	var start fs.FS
	if starterTemplatePath == "" {
		start, err = fs.Sub(starter, "starter/basic")
		mustNot(err)
	} else if starterTemplatePath[0] == ':' {
		start, err = fs.Sub(starter, "starter/"+starterTemplatePath[1:])
		mustNot(err)
	} else {
		start = os.DirFS(starterTemplatePath)
		stat, err := fs.Stat(start, ".")
		mustNot(err)
		if !stat.IsDir() {
			mustNot(fmt.Errorf("starter template must be a directory, a file was found at %s", starterTemplatePath))
		}
	}
	files, err := generate.PageScaffold(sitePkg, pageName, scaffold, start)
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
	fmt.Println("Found pages:", pageList)
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

// listPagesCmd prints the list of page names to stdout
func listPagesCmd(sitePkgRel string) {
	sitePkg, err := generate.GoListPackage(sitePkgRel)
	mustNot(err)
	pageList, err := generate.ScanSitemap(sitePkg)
	mustNot(err)
	fmt.Printf("%s", strings.Join(pageList, "\n"))
}

// deletePageCmd unlinks the page from the FS and updates the site pages
func deletePageCmd(pagePackage string) {
	pagePkg, err := generate.GoListPackage(pagePackage)
	mustNot(err)
	_, err = os.Stat(path.Join(pagePkg.Dir, "routemap.toml"))
	if os.IsNotExist(err) {
		log.Fatalf("Not a scaffold page '%s'", pagePackage)
	} else {
		mustNot(err)
	}
	var sitePkg generate.GoPackage
	{
		parts := strings.Split(pagePkg.ImportPath, "/")
		if len(parts) < 2 || parts[len(parts)-2] != "page" {
			mustNot(fmt.Errorf("invalid page path '%s'", pagePackage))
		}
		sitePkg, err = generate.GoListPackage(strings.Join(parts[:len(parts)-2], "/"))
		mustNot(err)
	}
	{
		// ask user if they want to delete the page folder
		fmt.Println("Found scaffold page at path ", pagePkg.Dir)
		fmt.Printf(">> Are you sure you want to delete this directory [yY]: ")
		var input string
		_, scErr := fmt.Scanln(&input)
		if scErr != nil || (input != "y" && input != "Y") {
			log.Fatalf("Cancelled delete, received answer '%s'\n", input)
		}
	}

	err = os.RemoveAll(pagePkg.Dir)
	mustNot(err)
	fmt.Println("Updating scaffold pages", sitePkg.ImportPath)
	pagesCmd(sitePkg.ImportPath)
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

// starterCmd generates a page strter template that can be used with the 'good page x --starter ...' command
func starterCmd(dest string) {
	start, err := fs.Sub(starter, "starter/default")
	mustNot(err)
	files, err := generate.StarterScaffold(dest, start)
	mustNot(err)
	err = generate.FlushFiles(".", files)
	mustNot(err)
	fmt.Println("Created files:")
	for _, file := range files {
		fmt.Println("\t-", filepath.Join(file.Dir, file.Name))
	}
	fmt.Println()

	fmt.Printf("Outputted default starter template to dest folder %s!", dest)
}

func mustNot(err error) {
	if err != nil {
		panic(err)
	}
}
