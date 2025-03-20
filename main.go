package main

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
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
	routes gen   <page_pkg>              enerate the routes.go file for a specified scaffold page package
	starter      <out_dir>               Create a new dir and populate it with a template for a custom starter page
	version                              Current version of the CLI

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
    -y
	Skip interactive confirmation message
    --no-resources
	Do not include a resources.go file in the new page
    --starter <dir path>
        Use specified directory as the starter template for the new page scaffold.

	There are a number of built-in page starters you can choose from. These are specified using
	a single colon prefix. Here is the list of available build-ins you can try:

		:basic                 No layout, just a simple page scaffold (default)
		:bootstrap5/layout     Useful Bootsrap v5.0 web console layout
		:bootstrap5/examples   Working demos with the Bootsrap v5.0 layout and components
		:bootstrap5/login      user login and registration flow with a mock mem-DB
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
	routesUsage = `usage: good routes gen <page_pkg_rel>

Generate golang code for the routemap.toml config file in a target page. This will also populate code
for any handlers or templates that are missing from the config.

Example
	good routes gen ./admin/site/page/example

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
	case "version":
		fmt.Printf("good version v0.1.2")

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
		var interactive bool
		if _, yesFlag := fArgs["-y"]; !yesFlag {
			interactive = generate.IsTTY()
		}

		starterTemplatePath := fArgs["--starter"]
		_, noResources := fArgs["--no-resources"]
		pageCmd(pArgs[1], pArgs[2], starterTemplatePath, !noResources, interactive)

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
		switch pArgs[1] {
		case "gen":
			if len(pArgs) < 3 {
				fmt.Println(routesUsage)
				log.Fatalln("Missing target page path")
			}
			routesGenCmd(pArgs[2])

		default:
			fmt.Println(routesUsage)
			log.Fatalf("Unknown routes command '%s'", pArgs[1])
		}

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
	// first try a go.mod in the CWD
	var (
		sitePkg generate.GoPackage
		err     error
	)
	if path, dir := generate.TryReadModFile(); path != "" {
		// use module info parsed from the go.mod file
		sitePkg, err = generate.ParseSitePackage(generate.GoModule{
			Path: path,
			Dir:  dir,
		}, sitePkgRel)
		userFail("reading project go.mod", err)
	} else {
		// use current package to find go module
		curPkg, err := generate.GoListPackage("./...")
		userFail("scanning subdirectires for a Go package", err)
		sitePkg, err = generate.ParseSitePackage(curPkg.Module, sitePkgRel)
		userFail("parsing module path", err)
	}
	err = generate.ValidateScaffoldLocation(sitePkg.Dir, scaffold)
	userFail("validating scaffold destination", err)
	files, err := generate.SiteScaffold(sitePkg, scaffold)
	userFail("generating scaffold files", err)
	pageFile, err := generate.PagesScaffold(sitePkg, []string{"intro"}, scaffold)
	mustNot(err)
	files = append(files, pageFile)
	start, err := fs.Sub(starter, "starter/intro")
	mustNot(err)
	pFiles, err := generate.PageScaffold(sitePkg, "intro", scaffold, start, false)
	mustNot(err)
	files = append(files, pFiles...)

	// FS operations
	err = generate.FlushFiles(sitePkg.Dir, files)
	userFail("writing scaffold files", err)

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
func pageCmd(sitePkgInput, pageName, starterTemplatePath string, withResources, interactive bool) {
	err := generate.ValidatePageName(pageName)
	userFail("validating your page name", err)
	sitePkg, err := generate.GoListPackage(sitePkgInput)
	userFail("loading the scaffold package detail", err)
	err = generate.ValidatePageLocation(filepath.Join(sitePkg.Dir, "page", pageName), scaffold)
	userFail("validating the page destination", err)

	if starterTemplatePath == "" {
		if interactive {
			// ask user what starter they would like to use
			fmt.Println("Choose a starter template")
			fmt.Println("Built-in options:")
			options := map[string]string{
				"1": ":basic                 No layout, just a simple page scaffold (default)",
				"2": ":bootstrap5/layout     Useful Bootsrap v5.0 web console layout",
				"3": ":bootstrap5/examples   Working demos with the Bootsrap v5.0 layout and components",
				"4": ":bootstrap5/login      user login and registration flow with a mock mem-DB",
				"5": ":minimum               Most bare bones option",
				"6": ":intro                 Introduction page to the good scaffold",
			}
			for i := 1; i <= len(options); i++ {
				fmt.Printf("\t[%d] %s\n", i, options[strconv.Itoa(i)])
			}
			fmt.Println("Select a built-in by number, blank for the default (:basic) or provide a page starter path")
			fmt.Printf("> ")
			input, err := generate.TimeoutScanln()
			if err != nil {
				if err.Error() == "unexpected newline" {
					// empty value, just fall back on default
					starterTemplatePath = ":basic"
				} else {
					userFail("reading your input", err)
				}
			} else if option, ok := options[input]; ok {
				starterTemplatePath = strings.SplitN(option, " ", 2)[0]
			} else {
				starterTemplatePath = input
			}
		} else {
			// just use the default for non-interactive terminals
			starterTemplatePath = ":basic"
		}
	}

	var start fs.FS
	if starterTemplatePath == "" {
		// this shouldn't happen
		mustNot(errors.New("empty starter template"))
	} else if starterTemplatePath[0] == ':' {
		start, err = fs.Sub(starter, "starter/"+starterTemplatePath[1:])
		userFail(fmt.Sprintf("loading internal starter template '%s'", starterTemplatePath), err)
	} else {
		start = os.DirFS(starterTemplatePath)
		stat, err := fs.Stat(start, ".")
		if err == nil && !stat.IsDir() {
			err = fmt.Errorf("starter template must be a directory, a file was found at path '%s'", starterTemplatePath)
		}
		userFail(fmt.Sprintf("reading custom starter directory path '%s'", starterTemplatePath), err)
	}

	if interactive {
		resourceMsg := "no"
		if withResources {
			resourceMsg = "yes"
		}
		// ask for confirmation
		fmt.Printf(
			strings.Join([]string{
				"Page Details",
				"\tscaffold            %q",
				"\tpage name           %q",
				"\tstarter template    %q",
				"\tresources.go        %s",
				"",
				"Create this page? [yY]: ",
			}, "\n"),
			sitePkg.ImportPath,
			pageName,
			starterTemplatePath,
			resourceMsg,
		)
		answer, err := generate.TimeoutScanln()
		userFail("reading your answer", err)
		if strings.ToUpper(answer) != "Y" {
			fmt.Println("Cancelled!")
			os.Exit(1)
			return
		}
	}

	files, err := generate.PageScaffold(sitePkg, pageName, scaffold, start, withResources)
	userFail("generating page scaffold", err)
	err = generate.FlushFiles(sitePkg.Dir, files)
	userFail("writing page files", err)

	pageImport := fmt.Sprintf("%s/page/%s", sitePkg.ImportPath, pageName)
	fmt.Printf("Created page at %s\n", pageImport)

	pageList, err := generate.ScanSitemap(sitePkg)
	mustNot(err)
	pages, err := generate.PagesScaffold(sitePkg, pageList, scaffold)
	userFail("generating new scaffold pages.go", err)
	err = generate.FlushFiles(sitePkg.Dir, []generate.File{pages})
	userFail("writing the updated pages.go file", err)

	stdout, err := generate.GoFormat(sitePkg.ImportPath + "/...")
	if err != nil {
		log.Fatalf("Page '%s' scaffold was create with fmt error: %s", pageImport, err)
	}
	if len(stdout) > 0 {
		fmt.Println("Output from go fmt:")
		fmt.Println(stdout)
	}
	fmt.Printf("Don't forget to run -> go generate %s/...\n", sitePkg.ImportPath)
	fmt.Printf("Created a good page for %s!", pageName)
}

// pagesCmd generates a new pages.go file by scanning the [site]/page/* directory
// for pages
func pagesCmd(sitePkgRel string) {
	sitePkg, err := generate.GoListPackage(sitePkgRel)
	userFail("scanning for a Go package at "+sitePkgRel, err)
	pageList, err := generate.ScanSitemap(sitePkg)
	userFail("scanning the scaffold for pages", err)
	fmt.Println("Found pages:", pageList)
	pages, err := generate.PagesScaffold(sitePkg, pageList, scaffold)
	userFail(fmt.Sprintf("generating updated pages.go for %v", pageList), err)
	err = generate.FlushFiles(sitePkg.Dir, []generate.File{pages})
	userFail("writing pages.go file", err)
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
	userFail(fmt.Sprintf("scanning for a Go package at '%s'", sitePkgRel), err)
	pageList, err := generate.ScanSitemap(sitePkg)
	userFail(fmt.Sprintf("scanning for pages in the scaffold '%s'", sitePkgRel), err)
	fmt.Printf("%s", strings.Join(pageList, "\n"))
}

// deletePageCmd unlinks the page from the FS and updates the site pages
func deletePageCmd(pagePackage string) {
	pagePkg, err := generate.GoListPackage(pagePackage)
	userFail(fmt.Sprintf("scanning for a Go package at '%s'", pagePackage), err)
	_, err = os.Stat(path.Join(pagePkg.Dir, "routemap.toml"))
	if os.IsNotExist(err) {
		err = fmt.Errorf("not a scaffold page '%s'", pagePackage)
	}
	userFail(fmt.Sprintf("reading scaffold page package '%s'", pagePackage), err)
	var sitePkg generate.GoPackage
	{
		parts := strings.Split(pagePkg.ImportPath, "/")
		if len(parts) < 2 || parts[len(parts)-2] != "page" {
			mustNot(fmt.Errorf("invalid page path '%s'", pagePackage))
		}
		sitePkg, err = generate.GoListPackage(strings.Join(parts[:len(parts)-2], "/"))
		userFail("loading details of the scaffold Go package", err)
	}
	if generate.IsTTY() {
		// ask user if they want to delete the page folder
		fmt.Println("Found scaffold page at path ", pagePkg.Dir)
		fmt.Printf(">> Are you sure you want to delete this directory [yY]: ")
		answer, err := generate.TimeoutScanln()
		userFail("reading your answer", err)
		if answer != "y" && answer != "Y" {
			log.Fatalf("Cancelled delete, received answer '%s'\n", answer)
		}
	}

	err = os.RemoveAll(pagePkg.Dir)
	userFail(fmt.Sprintf("removing page directory '%s'", pagePkg.Dir), err)
	fmt.Println("Updating scaffold pages", sitePkg.ImportPath)
	pagesCmd(sitePkg.ImportPath)
}

// routesCmd will parse a routemap.toml file and generate routes, handlers and templates
// as needed
func routesGenCmd(pagePkgRel string) {
	pkg, err := generate.GoListPackage(pagePkgRel)
	userFail(fmt.Sprintf("scanning for a Go package at '%s'", pagePkgRel), err)
	routemapPath := filepath.Join(pkg.Dir, "routemap.toml")
	routemapContent, err := os.ReadFile(routemapPath)
	userFail(
		fmt.Sprintf("reading routemap file at path '%s'", routemapPath),
		err,
	)
	tree, err := toml.LoadBytes(routemapContent)
	userFail(
		fmt.Sprintf("parsing TOML format of file '%s'", routemapPath),
		err,
	)
	pageName := pkg.Name
	sitePkg, err := generate.SiteFromPagePackage(pkg)
	userFail(
		fmt.Sprintf("loading the scaffold Go package for page '%s'", pkg.ImportPath),
		err,
	)
	_, rscErr := os.Stat(filepath.Join(pkg.Dir, "resources.go"))
	hasResources := !os.IsNotExist(rscErr)
	pageRoutes, missTpl, missHlr, err := routemap.ProcessRoutemap(
		tree,
		filepath.Join("page", pageName, "templates"),
		hasResources,
	)
	userFail(
		fmt.Sprintf("parsing routemap views for file '%s'", routemapPath),
		err,
	)
	entries, routes, handlers, templates, err := routemap.TemplateDataForRoutes(pageRoutes, missTpl, missHlr)
	userFail(
		fmt.Sprintf("processing routemap details for file '%s'", routemapPath),
		err,
	)
	files, err := generate.RoutesScaffold(
		sitePkg,
		pageName,
		entries,
		routes,
		handlers,
		templates,
		scaffold,
		hasResources,
	)
	userFail(fmt.Sprintf("generating routemap files for page '%s'", pkg.ImportPath), err)
	if len(missHlr)+len(missTpl) > 0 {
		modifiedContent, err := routemap.ModifiedRoutemap(bytes.NewReader(routemapContent), missTpl, missHlr)
		userFail(
			fmt.Sprintf("modifying routemap contents for file '%s'", routemapPath),
			err,
		)
		files = append(files, generate.File{
			Dir:       filepath.Join("page", pageName),
			Name:      "routemap.toml",
			Contents:  modifiedContent,
			Overwrite: true,
		})
	}

	// write files to disk
	err = generate.FlushFiles(sitePkg.Dir, files)
	userFail(fmt.Sprintf("Writing routing related files for page '%s'", pkg.Dir), err)
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
	start, err := fs.Sub(starter, "starter/basic")
	mustNot(err)
	files, err := generate.StarterScaffold(dest, scaffold, start)
	userFail("generating starter template files to '"+dest+"'", err)
	err = generate.FlushFiles(".", files)
	userFail("writing starter template files to '"+dest+"'", err)
	fmt.Println("Created files:")
	for _, file := range files {
		fmt.Println("\t-", filepath.Join(file.Dir, file.Name))
	}
	fmt.Println()

	fmt.Printf("Outputted default starter template to dest folder %s!", dest)
}

// mustNot will panic if the error is not nil. The panic will produce a stack trace
// of all running Goroutines, hence it should be used when a non-nil error means there
// is likely a bug in the program that is not likely to be meaningful to the CLI user.
func mustNot(err error) {
	if err != nil {
		panic(err)
	}
}

// userFail will log an error for the user to see, describing what task was being attempted
// when the error was encountered. This will not print a stack trace. It should be used when
// an error might be expected and the user *might* be able to do something about it.
func userFail(task string, err error) {
	if err == nil {
		return
	}
	log.Fatalf("Failed while %s. Error: %s", task, err)
}
