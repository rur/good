[![Build Status](https://travis-ci.com/rur/good.svg?token=ghq4t9FLdVA8tqkRUMoY&branch=main)](https://travis-ci.com/rur/good)

# <img src="docs/check.svg" height="23rem" width="30rem"/>Good

## A pretty good web scaffold for Golang

CLI tool for embedding a web setup in a Golang application.

- Generates plain, grok-able code
- Basic dependencies
- Embedded server assets
- Classic template composition
- No surprises

_[see [developer notes](#developer-notes)]_

### Overview

The `good scaffold` command outputs files for a web server to a sub package of a
Go project. It geared towards under-the-hood GUIs like admin tools, consoles and user dashboards.

Top-level pages are added to the site using the `good page` command, which has a
[starter template](#Starter-Template) feature to speed-up development. The scaffold
is designed to be integrated with the codebase of a typical web service or background daemon
implemented in Go.

#### TLDR; quickstart

    $ go get github.com/rur/good
    $ cd ~/path/to/mygoproject
    [mygoproject]$ good scaffold ./myportal
    [mygoproject]$ go generate ./myportal/...
    [mygoproject]$ go run ./myportal

Visit localhost:8000 and take it from there.

### CLI Overview

The CLI tools will generate HTML and Golang files that you should modify & refactor
to suit your needs.

> Tip: _Use the `-h` flag for help with commands_

#### Good Scaffold ...

    $ good scaffold ./portal

Example command that creates a new app at `[current_go_mod]/portal`.

#### Good Page ...

    $ good page ./portal settings

Add a new 'settings' page to an existing scaffold path.

#### Good Routes ...

    $ good routes ./portal/page/settings

Re-generate the routing code for the portal settings page based on the
`./portal/page/settings/routemap.toml` file.

#### Good Starter ...

    $ good starter ./portal/my-page-starter

Make a directory and populate it will code template files for a customized starter page
that can be used with the 'good page' command like so.

    $ good page ./portal mynewpage --starter-template ./portal/my-page-starter

## Developer Notes

#### 1. Generate plain, grok-able code

The output is mostly vanilla Golang and HTML templates. We embrace some redundancy
so that the code will be more static, easier to read and customize.
This works well with the Golang type system and tooling, which makes refactoring a cinch
at the cost of some extra typing.

#### 2. Basic Dependencies

We avoid a lot of mandatory dependencies by taking full advantage of the standard library.
There is no plugin system; if you are familiar with Golang you can rely on an easy-to-follow
codebase to integrate your chosen libraries manually.

#### 3. Binary Embedded

The `//go:embed ` compiler directive is configured so that web server assets can be fully embedded at compile time.
This gives you the option to distribute your GUI as a self-contained binary.

#### 4. Classic template composition

HTML template composition has excellent support in Golang. This scaffold uses the
[Treetop library](https://github.com/rur/treetop) to help organize templates, with
the added benefit of fragment hot-swapping to enhance interactivity.

#### 5. No Surprises

This scaffold is more of a workhorse than a unicorn; we embrace some practical
limitations for the purpose of simplicity and embedding.
You should take care to judge whether this will be a good fit for your project.

## Starter Template

[TODO] Docs
