[![Build Status](https://travis-ci.com/rur/good.svg?token=ghq4t9FLdVA8tqkRUMoY&branch=main)](https://travis-ci.com/rur/good)

# <img src="docs/readme_logo.svg" alt="Good Web Scaffold"/>

## A pretty good web scaffold for Golang

Tools for embedding a web GUI in a Golang project.

- CLI that generates plain, grok-able code
- Self contained server with basic dependencies
- Embedded assets
- Classic HTML template composition
- Pretty boring, avoid surprises

[ [CLI Overview](#cli-overview) ~
[Building Pages](#building-pages) ~
[Developer Notes](#developer-notes) ~
[Starter Templates](#starter-templates)]

### Overview

The `good scaffold` command outputs files for a web server to a sub package of a
Go module. The GUI scaffold is geared towards Golang web services
or daemons, for building web consoles and admin tools.

Top-level pages are added using the `good page` command, which has a
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

#### Good Scaffold \<site\>

The snippet creates a new app at `./portal`, relative to the current Go module.

    $ good scaffold ./portal

#### Good Page \<site\> \<name\>

Adds a new 'settings' page to our `./portal` example.

    $ good page ./portal settings --starter :bootstrap5/layout

#### Good Pages {gen list delete}

Utilities for site pages. For example, list the pages to stdout

    $ good pages list ./portal
    home
    settings

#### Good Routes \<page\>

Re-generate the route plumbing code for our `./portal` settings page. The tool will read the `./portal/page/settings/routemap.toml` config file.

    $ good routes ./portal/page/settings

#### Good Starter \<outdir\>

    $ good starter ./portal/my-page-starter

Create a directory and populate it with code template files for a customized starter page.
This can be used with the Good Page command like so.

    $ good page ./portal mynewpage --starter ./portal/my-page-starter

## Building Pages
> We plan to add a web console in v2
>
>_– Opimistic BE Developer_

With this scaffolding tool, we hope to make GUI development more enjoyable and
approachable for Backend devs.

* Quick-start boilerplate for dashboard layouts:
  * Get set up with a suitable CSS framework
  * Functioning examples, share starter templates
* Code generation with [TOML](https://toml.io/en/) config:
  * Make full use of the compiler
  * Unify config for routing, templates and dispatch
* Minimize logic outside of Golang
  * Server-side rendering
  * Use existing project code directly

### TOML Routemaps
Routemaps combine code generation with a familiar approach to HTML templating.
The `good routes` command reads this config and generates plumbing code.

This sample has two endpoints; dive into the scaffold to see more.

```TOML
_ref = "base"
_template = "/templates/base.html.tmpl"
_handler = "baseHandler"

  [content]
  _ref = "main-content"
  _template = "/templates/content/main.html.tmpl"
  _handler = "mainContentHandler"
  _path = "/example"

  [content]
  _ref = "other-content"
  _template = "/templates/content/other.html.tmpl"
  _handler = "otherContentHandler"
  _path = "/other-example"
```

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
This gives you the option to distribute your GUI as a self-contained binary.

#### 4. Classic HTML template composition

HTML template composition has excellent support in Golang. Our scaffold uses the
[Treetop library](https://github.com/rur/treetop) to help organize templates and handlers,
with the added benefit of fragment hot-swapping for enhanced interactivity.

#### 5. Pretty boring, avoid surprises

The scaffold is more of a workhorse than a unicorn; we embrace some practical
limitations for the purpose of tight server-side integration.
Take care to judge the limitations for yourself and decide what is right for your project.

## Starter Templates

[TODO] Docs
