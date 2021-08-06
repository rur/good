[![Build Status](https://travis-ci.com/rur/good.svg?token=ghq4t9FLdVA8tqkRUMoY&branch=main)](https://travis-ci.com/rur/good)

# <img src="docs/readme_logo.svg" alt="Good Web Scaffold"/>

## A pretty good web scaffold for Golang

Tools for embedding a web GUI in a Golang project.

- CLI that generates plain, grok-able code
- Self contained server with basic dependencies
- Embedded assets
- Classic HTML template composition
- Pretty boring, no surprises

(see [developer notes](#developer-notes))

### Overview

The `good scaffold` command outputs files for a web server to a sub package of a
Go module. The scaffold is geared towards user consoles and admin tools; integrated
GUIs for a Golang web service or daemon.

Top-level pages are added to the site using the `good page` command, which has a
[starter template](#Starter-Template) feature to speed-up development. The scaffold package
is unobtrusive so that it can be embedded in the codebase of an existing system.

#### TLDR; quickstart

    $ go get github.com/rur/good
    $ cd ~/path/to/mygoproject
    [mygoproject]$ good scaffold ./myportal
    [mygoproject]$ go generate ./myportal/...
    [mygoproject]$ go run ./myportal --port 8000

Visit localhost:8000 and take it from there.

### CLI Overview

The CLI tools will generate HTML and Golang files that you should modify & refactor
to suit your needs.

> Tip: _Use the `-h` flag for help with commands_

#### Good Scaffold \<site\>

Create a new app at `./portal` relative to the current project.

    $ good scaffold ./portal

#### Good Page \<site\> \<name\>

Add a new 'settings' page to an existing scaffold path.

    $ good page ./portal settings

#### Good Pages {gen list delete}

Utilities for site pages. For example, list the pages to stdout

    $ good pages list ./portal
    home
    settings

#### Good Routes \<page\>

Re-generate the route plumbing code for the portal settings page based on the
`./portal/page/settings/routemap.toml` file.

    $ good routes ./portal/page/settings

#### Good Starter \<outdir\>

    $ good starter ./portal/my-page-starter

Make a directory and populate it with code templates for a customized starter page
that can be used with the 'good page' command like so.

    $ good page ./portal mynewpage --starter-template ./portal/my-page-starter

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

HTML template composition has excellent support in Golang. The scaffold uses the
[Treetop library](https://github.com/rur/treetop) to help organize templates and handlers,
with the added benefit of fragment hot-swapping for enhanced interactivity.

#### 5. Pretty boring, no surprises

This scaffold is more of a workhorse than a unicorn; we embrace some practical
limitations for the purpose of tight server-side integration.
Take care to judge the limitations for yourself and decide what is right for your project.

## Starter Template

[TODO] Docs
