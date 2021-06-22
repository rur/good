[![Build Status](https://travis-ci.com/rur/good.svg?token=ghq4t9FLdVA8tqkRUMoY&branch=main)](https://travis-ci.com/rur/good)

# _Good_

## A pretty good web scaffold for Golang

- Clear, grok-able code
- Basic dependencies
- Embed within your binary
- Classic template composition
- No surprises

`good` is a code-gen tool for embedding a web GUI in an existing Golang application.
It outputs a self-contained scaffold for a modern server side app, with the aim of being
straightforward to setup and manage over time.

The general setup is geared towards: service consoles, admin workflows
and user controls with direct server integration.

### CLI Overview

Commands that generate HTML and Golang files to be modified & refactored to suit your needs.

#### TLDR; quickstart

    $ go get github.com/rur/good
    $ cd ~/path/to/mygoproject
    [mygoproject]$ good scaffold ./myportal
    [mygoproject]$ go generate ./myportal/...
    [mygoproject]$ go run ./myportal

Visit localhost:8000 and take it from there.

#### Good Scaffold ...

    $ good scaffold ./portal dashboard

Create a new app at `[current_go_mod]/portal` with a single page named _dashboard_.

#### Good Page ...

    $ good page ./portal settings

Add a new 'settings' page to the existing portal scaffold.

#### Good Routes ...

    $ good routes ./portal/settings

Re-generate the routing code for the portal settings page based on the
`./portal/settings/routemap.toml` file.

## Folder Structure

### Routemap TOML

Routing is b

## Intro

#### 1. Generate obvious, grok-able code

The output is mostly vanilla Golang and HTML. We embrace a little redundancy
so that the code will be more static, easier to customize and more obvious. This works
well with the Golang tooling, which makes refactoring a cinch.

#### 2. Only Basic Dependencies

It is simpler to add the libraries you need, if there are fewer there to begin with.
We use a library for managing handlers ([treetop](https://github.com/rur/treetop)), and the go
standard library for the rest.

#### 3. Binary Embedded

The `//go:embed ` directive is configured so that the web server is fully embedded at compile time.

#### 4. Classic template composition

Only one way in or out; uniform endpoints greatly reduce the surface area you need to think about.

#### 5. No Surprises

If a framework might be overkill for your project, this bootstrap is a pretty good alternative.

## QnA

#### Is this a framework or a static site generator?

Neither, it's just a scaffold. Once the code is generated it's yours, adapt it to your needs.
Routing is the only untouchable code, mostly because it is tedious and not very interesting.

#### Why use a special library to manage handlers?

Basically, I don't want to live without plumbing! The coordination of handlers and templates with
request routing needs to be declarative (otherwise 🤯). [Treetop](https://github.com/rur/treetop)
does this without transitive dependencies.

#### Where does the frontend go?

TODO: ...

#### The generator creates lots of handlers, why so much bloat?

I take the view that DRY code ariases from repeated refactoring (see TDD).
Generating a bunch of simple working code is a nice way to start the refactoring process,
but the rest is up to the programmer.
