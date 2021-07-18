package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
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
	scaffold <site_pkg> [<page_name>...]   Create a new site scaffold at at a package relative to the working dir
	page     <site_pkg> <page_name>        Add a new page to an existing scaffold
	pages    <site_pkg>                    Generate a routes.go file from a TOML config
	routes   <page_pkg>                    Generate a routes.go file from a TOML config
	starter  <out_dir>                     Create a new dir and populate it with a template for a custom starter page

`
	scaffoldUsage = `usage: good scaffold <site_pkg_rel> [<page_name>...]

Create a new site scaffold in the current Golang project

Example
	good scaffold ./admin/site dashboard settings

Arguments
	site_pkg_rel    relative import path of the new site from the current Go module
	page_name       optional list of page names to be initialized along with the site, default is 'example'

`
	pageUsage = `usage: good page <site_pkg_rel> <page_name> [--starter-template <path>]

Add a new page to an existing scaffold site.

Example
	good page ./admin/site settings

Arguments
	site_pkg_rel   relative import path of an existing scaffold site from the current Go module
	page_name      package name of the new page to initialize

Options
    --starter-template <path>
        Use specified directory as the starter template for the new page scaffold

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
	starterUsage = `usage: good starter <out_dir>

Generate a template directory for a custom starter page that can be used with the
'goodo page' comand.


Example
	$ good starter ./admin/utils/customstarter
    [...customize the files...]

    # create a new page using your starter
    $ good page ./admin/site awesomenewpage --starter-template ./admin/utils/customstarter

Arguments
	out_dir   a not-already-existing path where a folder will be created

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
			for j := i; j < len(pArgs)-1; j += 2 {
				fArgs[pArgs[j]] = pArgs[j+1]
			}
			pArgs = pArgs[:i]
			break
		}
	}

	if len(pArgs) < 1 {
		fmt.Println(usage)
		log.Fatalf("Missing <command>")
	}
	switch pArgs[0] {
	case "scaffold":
		if len(pArgs) < 2 {
			fmt.Println(scaffoldUsage)
			log.Fatalf("Missing target site package path")
		}
		scaffoldCmd(pArgs[1], pArgs[2:])

	case "page":
		if len(pArgs) < 3 {
			fmt.Println(pageUsage)
			log.Fatalf("Missing required arguments")
		}
		starterTemplatePath := fArgs["--starter-template"]
		fmt.Printf("Starter template %#v \n", starterTemplatePath)
		pageCmd(pArgs[1], pArgs[2], starterTemplatePath)

	case "pages":
		if len(pArgs) < 2 {
			fmt.Println(pagesUsage)
			log.Fatalf("Missing target site package path")
		}
		pagesCmd(pArgs[1])

	case "routes":
		if len(pArgs) < 2 {
			fmt.Println(routesUsage)
			log.Fatalln("Missing target scaffold page path")
		}
		routesCmd(pArgs[1])

	case "starter":
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
func scaffoldCmd(sitePkgRel string, pages []string) {
	if len(pages) == 0 {
		// if no page names were listed, add a page called 'example'
		pages = []string{"example"}
	}

	// use current package to find go module
	curPkg, err := generate.GoListPackage(".")
	mustNot(err)
	sitePkg, err := generate.ParseSitePackage(curPkg.Module, sitePkgRel)
	mustNot(err)
	err = generate.ValidateScaffoldLocation(sitePkg.Dir, scaffold)
	mustNot(err)
	files, err := generate.SiteScaffold(sitePkg, scaffold)
	mustNot(err)
	pageFile, err := generate.PagesScaffold(sitePkg, pages, scaffold)
	mustNot(err)
	files = append(files, pageFile)

	start, err := fs.Sub(starter, "starter/default")
	mustNot(err)

	for _, page := range pages {
		err = generate.ValidatePageName(page)
		mustNot(err)
		pFiles, err := generate.PageScaffold(sitePkg, page, scaffold, start)
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
func pageCmd(sitePkgRel, pageName, starterTemplatePath string) {
	err := generate.ValidatePageName(pageName)
	mustNot(err)
	sitePkg, err := generate.GoListPackage(sitePkgRel)
	mustNot(err)
	err = generate.ValidatePageLocation(filepath.Join(sitePkg.Dir, "page", pageName), scaffold)
	mustNot(err)
	var start fs.FS
	if starterTemplatePath != "" {
		start = os.DirFS(starterTemplatePath)
		stat, err := fs.Stat(start, ".")
		mustNot(err)
		if !stat.IsDir() {
			mustNot(fmt.Errorf("starter-template must be a directory, a file was found at %s", starterTemplatePath))
		}
	} else {
		start, err = fs.Sub(starter, "starter/default")
		mustNot(err)
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

// starterCmd generates a page strter template that can be used with the 'good page x --starter-template ...' command
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
