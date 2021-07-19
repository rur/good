[![Build Status](https://travis-ci.com/rur/good.svg?token=ghq4t9FLdVA8tqkRUMoY&branch=main)](https://travis-ci.com/rur/good)

# _Good_

## A pretty good web scaffold for Golang

`good` is a CLI tool for creating an embedded web console for a Golang application.

- Generates plain, grok-able code
- Only basic dependencies
- Embed all assets
- Classic template composition
- No surprises

_[see [developer notes](#notes)]_

### Overview

The `good scaffold` command outputs files for a web server in a sub package of a
Go project. Pages are added to the site using the `good page` command, which has a
starter template feature to help speed-up development.

The scaffold is a general purpose setup that is particularly useful for apps that have
under-the-hood integrations like admin tools and service dashboards.

#### TLDR; quickstart

    $ go get github.com/rur/good
    $ cd ~/path/to/mygoproject
    [mygoproject]$ good scaffold ./myportal
    [mygoproject]$ go generate ./myportal/...
    [mygoproject]$ go run ./myportal

Visit localhost:8000 and take it from there.

### CLI Commands

The CLI tools will generate HTML and Golang files that you should modify & refactor
to suit your needs.

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

## Developer Notes

#### 1. Generate plain, grok-able code

The output is mostly vanilla Golang and HTML. We embrace some redundancy
so that so that the code will be more static, easier to read and customize.
This works well with the Golang type system and tooling, making refactoring a cinch.

#### 2. Basic Dependencies

Taking advantage of the standard library helps us to avoid a lot of mandatory dependencies.
There is no plugin system, instead we encourage you to manually integrate your chosen libraries
and rely on the easy-to-follow codebase to help you out.

#### 3. Binary Embedded

The `//go:embed ` directive is configured so that web server assets can be fully embedded at compile time.

#### 4. Classic template composition

HTML template composition has excellent support in Golang. This scaffold uses the
[Treetop library](https://github.com/rur/treetop) to help organize templates, with
the added benefit of fragment hot-swapping to enhance interactivity.

#### 5. No Surprises

This scaffold is more of a workhorse than a unicorn, we embrace many practical
limitations for the benefit of long term maintenance and tight integration.
Take care to judge whether this will be a good fit for your project or not.
