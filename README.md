[![Build Status](https://app.travis-ci.com/rur/good.svg?branch=main)](https://app.travis-ci.com/rur/good)

# <img src="docs/readme_logo.svg" alt="Good Web Scaffold"/>

## A pretty good web scaffold for Golang

The `good` tool embeds a web GUI in a Golang project.

- CLI that generates plain, grok-able code
- Self contained with basic dependencies
- Fully embedded assets and templates
- Classic HTML template composition
- Boring tech, fewer surprises

[ [CLI Overview](#cli-overview) ~
[Building GUIs](#building-guis) ~
[Starter Templates](#starter-templates) ~
[Developer Notes](#developer-notes) ]

### Overview

The `good scaffold` command outputs a web GUI setup to a sub package of a
Go project. The scaffold is geared towards web consoles and admin tools
for backend services or daemons.

Pages are added using the `good page` command, which has a
[starter template](#starter-templates) feature to speed-up development. The scaffold package
is unobtrusive so that it can be embedded in an existing codebase without difficulty.

#### TLDR; quickstart

    $ go get github.com/rur/good
    $ cd ~/path/to/mygoproject
    [mygoproject]$ good scaffold ./myportal
    [mygoproject]$ go generate ./myportal/...
    [mygoproject]$ go run ./myportal --dev --port 8000

Visit `localhost:8000` and take it from there.

## CLI Overview

The CLI tools will generate HTML and Golang files that you can modify & refactor
to suit your needs.

> Tip: _Use the `-h` flag for help with commands_

#### Good Scaffold \<scaffoldPkg\>

This snippet creates a new scaffold at `./portal`, relative to the current Go module.

    $ good scaffold ./portal

#### Good Page \<scaffoldPkg\> \<pagename\>

Adds a new 'settings' page to our `./portal` example.

    $ good page ./portal settings

#### Good Pages {gen list delete} ...

Utilities for site pages. For example, list the pages to stdout

    $ good pages list ./portal
    home
    settings

#### Good Routes Gen \<pagePkg\>

Update the routing code for the specified scaffold page package.

Generate a __routes.go__ file from the `routemap.toml` config file in the target page. This
contains layout and endpoint plumbing code that uses the [Treetop](https://github.com/rur/treetop) library.


    $ good routes gen ./portal/page/settings

#### Good Starter \<outdir\>

    $ good starter ./portal/my-page-starter

Create a directory and populate it with code template files for a customized starter page.
This can be used with the Good Page command like so.

    $ good page ./portal mynewpage --starter ./portal/my-page-starter

## Building GUIs

> We'll add a web console in the _next_ version
>
>_– Backend Developer_

This scaffold reduces the overhead of adding a GUI to your
project in the following ways:

* Minimize the code outside of Golang,
  * Use more of your available project code
  * Make it convenient to do templating on the server-side
* Quick-start page templates:
  * Pick a suitable layout and CSS toolkit
  * Share examples for common workflows
* Code generation with TOML config:
  * Static plumbing code for page routes 
  * Take advantage of the compiler

#### Low Cost Interactivity

The server is able to render fragments of your page layout to be applied to the browser DOM.
It is possible to satisfy most UX requirements this way, without the need to ship an API for your GUI
or the corresponding IO boilerplate. 

Try the [example starter page](starter#bootstrap-v5)

### Treetop Routemap Layouts

Routemaps combine code generation with the classic _hierarchical_ approach to web
templating. The [Treetop library](https://github.com/rur/treetop) is used to build layouts and handlers
for endpoints.

The command `good routes gen` will read a TOML config for a page and generate endpoint plumbing code.

#### Routemap preview

```TOML
_ref = "base"
_template = "/templates/base.html.tmpl"
_handler = "baseHandler"

  [[content]]
  _ref = "main-content"
  _template = "/templates/content/main.html.tmpl"
  _handler = "mainContentHandler"
  _path = "/example"

  [[content]]
  _ref = "other-content"
  _template = "/templates/content/other.html.tmpl"
  _handler = "otherContentHandler"
  _path = "/other-example"
```

This basic config has two endpoints that share the same base view. Templates can be nested further.
[Try the scaffold](#tldr-quickstart) and explore some working examples.


## Starter Templates

Page boilerplate can be loaded from a local folder or using one of the built-in options.


### Built-in page starter

See the [starter#good-starter](starter/README.md) for details about what built-in options are available.

```
good page ./mysite mynewpage --starter :bootstrap5/layout
```

### Custom page starter

Commit a starter page to your project with custom boilerplate. Very useful if
you like to do a lot of prototyping!

The `good starter` command will help you to get set up.


## Developer Notes

#### 1. CLI that generates plain, grok-able code

The output is mostly vanilla Golang and HTML templates. We embrace some redundancy
so that the code will be more static, easier to read and customize.
This works well with the Golang type system and tooling, which makes refactoring a cinch
at the cost of some extra typing.

#### 2. Self-contained server with basic dependencies

We avoid a lot of mandatory dependencies by taking full advantage of the standard library.
There is no plugin system; if you are familiar with Golang you can rely on an easy-to-follow
codebase to integrate your chosen libraries manually.

#### 3. Embedded assets

The `//go:embed ` compiler directive <sup>[≥[go1.16](https://golang.org/doc/go1.16#library-embed)]</sup>
is configured so that web server assets can be fully embedded at compile time.
This gives you the option to distribute your GUI as a self-contained binary or embedded in an
another Go server.

#### 4. Classic HTML template composition

HTML template composition has excellent support in Golang. This scaffold uses the
[Treetop library](https://github.com/rur/treetop) to help organize templates and handlers,
with the added benefit of fragment hot-swapping for enhanced interactivity.

#### 5. Boring tech, fewer surprises

The scaffold is more of a workhorse than a unicorn; we embrace some practical
limitations for the purpose of tight server-side integration.
Take care to judge the limitations for yourself and decide what is right for your project.
