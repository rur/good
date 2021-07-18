[![Build Status](https://travis-ci.com/rur/good.svg?token=ghq4t9FLdVA8tqkRUMoY&branch=main)](https://travis-ci.com/rur/good)

# _Good_

## A pretty good web scaffold for Golang

`good` is a CLI tool for building an embedded web console in a Go application.

- Generates plain, grok-able code
- Only basic dependencies
- Embed all assets
- Classic template composition
- No surprises

The `good scaffold` command outputs files for a web server in the sub package of a
Go project. Pages are added to the site using the `good page` command, which has a
starter template feature to help speed-up development.

The scaffold is a general purpose web setup that is particularly useful for apps with
under-the-hood integrations like admin tools and service consoles.

### CLI Overview

The CLI tool will generate HTML and Golang files in your project,
which you should modify & refactor to suit your needs.

#### TLDR; quickstart

    $ go get github.com/rur/good
    $ cd ~/path/to/mygoproject
    [mygoproject]$ good scaffold ./myportal
    [mygoproject]$ go generate ./myportal/...
    [mygoproject]$ go run ./myportal

Visit localhost:8000 and take it from there.

#### Good Scaffold ...

    $ good scaffold ./portal

Create a new app at `[current_go_mod]/portal`.

#### Good Page ...

    $ good page ./portal settings

Add a new 'settings' page to an existing scaffold path.

#### Good Routes ...

    $ good routes ./portal/page/settings

Re-generate the routing code for the portal settings page based on the
`./portal/settings/routemap.toml` file.

#### Good Starter ...

    $ good starter ./portal/my-page-starter

Create a directory containing files for a custom starter page that can be used with the
'good page' command like so.

    $ good page ./portal mypage --starter-template ./portal/my-page-starter

### Intro

#### 1. Generate clear, grok-able code

The output is mostly vanilla Golang and HTML. We embrace a little redundancy
so that the code will be more static and so easier to customize.
This works very well with the Golang tooling, making refactoring a cinch.

#### 2. Basic Dependencies

We take full advantage of the standard _html/template_ and _net/http_
libraries to avoid many dependencies. The scaffold code is clearly commented
to make it easy for you to integrate the dependencies you wish to use in
your project.

#### 3. Binary Embedded

The `//go:embed ` directive is configured so that the web server is fully embedded at compile time.

#### 4. Classic template composition

Nested HTML template composition is a tried a true approach for building
web GUI that has excellent support in Golang. With the addition of fragment
hot-swapping, nested templates are capable of delivering a modern web experience
from a server-side app.

#### 5. No Surprises

If a framework might be overkill for your project, and you need more than static
content pages, this scaffold is probably a reasonable option.
