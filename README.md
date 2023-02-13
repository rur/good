[![Build Status](https://rur.semaphoreci.com/badges/good/branches/master.svg?key=12cc127d-6aa0-4cae-8e0b-d667acfcdb4b)](https://rur.semaphoreci.com/projects/good)

# <img src="docs/readme_logo.svg" alt="Good Web Scaffold"/>

## A good-enough web scaffold for Golang web consoles

The `good` tool embeds a web GUI in an existing Golang project.

- Self-contained with minimal dependencies
- Plain, grok-able code
- Fully embedded assets and templates
- Classic HTML template composition
- Boring tech, few surprises

[ [CLI Overview](#cli-overview) ~
[Building GUIs](#building-guis) ~
[Starter Templates](#starter-templates) ~
[Developer Notes](#developer-notes) ]

### Overview

The `good scaffold` command outputs a web GUI setup to a sub package of a
Go project.

The scaffold package
is unobtrusive so that it can be embedded in an existing Go module and make use of the project
boilerplate.

Pages are added using the `good page` command, which has a
[starter template](#starter-templates) feature to speed things up.

#### TLDR; quickstart

    $ go install github.com/rur/good
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

Generate a **routes.go** file from the `routemap.toml` config file in the target page. This
contains layout and endpoint plumbing code that uses the [Treetop](https://github.com/rur/treetop) library.

    $ good routes gen ./portal/page/settings

#### Good Starter \<outdir\>

    $ good starter ./portal/my-page-starter

Create a directory and populate it with code template files for a customized starter page.
This can be used with the Good Page command like so.

    $ good page ./portal mynewpage --starter ./portal/my-page-starter

## Building GUIs

> We'll add a console in the next version _– Go Developer_

This scaffold aims to reduce the overhead of adding a GUI to your
Go project in the following ways:

- Quick-start page templates,
  - Start with a workable layout and CSS toolkit.
- Code generation with TOML config,
  - Static plumbing code for page routes (you don't want to do this manually).
  - We try to balance complexity between code-gen, compile time and runtime.
- Keep the bulk of the business logic in Go,
  - Make it convenient to do GUI templating on the server-side.
  - Above all, avoid the multiplication of IO boilerplate.

#### Light weight interactivity

In addition to standard HTML page requests, the web server supports a custom protocol extension enabling
HTML fragments to be projected to the client. Many modern UX requirements can be satisfied this way,
alleviating the need for a more fully fledged interactive approach.

_Note:_ HTML fragments are an opt-in feature at the endpoint level, see `_partial` and `_fragment`
flags in the Routemap config guide.

#### Serving SPAs

If a Single Page App is what you have in mind, a simple web scaffold is a great way to serve your app container
templates, along with top-level nav and any auxiliary content, without incorporating a full-featured
HTML web framework with your codebase.

### Routemap Layouts

Routemaps combine code generation with the classic 'layout hierarchy' approach to web
templating. The [Treetop library](https://github.com/rur/treetop) is used to bind template files
and handlers to HTTP endpoints.

The `good routes gen`command uses a TOML config to generate the endpoint plumbing code for a page.

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

^ This basic example has two endpoints that share the same base view. Templates can be nested further.
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
so that the code will be more static, and easier to read and customize.
This works well with the Golang type system and tooling, which makes refactoring a cinch
at the cost of some extra typing.

#### 2. Self-contained server with minimal dependencies

We avoid a lot of mandatory dependencies by making the most of the standard library.
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

#### 5. Boring tech, few surprises

The scaffold is more of a workhorse than a unicorn; we embrace some practical
limitations for tight server-side integration.
Take care to judge the limitations for yourself and decide what is right for your project.
