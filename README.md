# _Good_

## A pretty good web scaffold for Golang

1. Generate obvious, grok-able code
1. Only basic dependencies
1. Embedded in the binary
1. Easy to secure
1. No surprises

`good` is a code-gen tool for adding a web GUI to a Golang project. 
The scaffold is geared towards integrated apps like admin portals
or service dashboards. It's a general purpose setup however, with a focus
on being low-maintenance over time.

### CLI Overview

These commands support a rapid development flow: generate code and refactor to suit your needs.

#### Good Scaffold

    $ good scaffold portal

Create a new app at `[pkg]/portal`. The namespace is read from `./go.mod`.

#### Good Page

    $ good page ./portal mypage

Add a new page 'mypage' to the portal (contains its own route config).

#### Good Routes

    $ good routes ./portal/mypage/routemap.toml

(Re)generate the routing code from a config file. This will overwrite `./portal/mypage/routemap.go`
and output any handler functions that are missing.

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

#### 4. Easy to Secure

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

