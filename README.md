# _Good_

## A pretty good web scaffold for Golang

1. Generates obvious, grok-able code
1. Only basic dependencies
1. Binary embedded
1. Easy to secure
1. No surprises

`good` is a code-gen tool for embedding a web GUI in your Golang project. The scaffold
is geared towards an admin portal or server dashboard. However, it is a general purpose setup,
useful when a framework might be over-doing it.

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

The output is mostly vanilla Golang and HTML. Code-gen helps us to embrace a little
redundancy so our code can be more static and easier to modify. Go tooling and builtin HTML templates
make it fun to generate and refactor new endpoints as needed. (It can also help performance ⚡️)

#### 2. Only Basic Dependencies

With few mandatory dependencies<sup>[1]</sup> it is more straightforward to
add the other libraries you might want.

#### 3. Binary Embedded

The `//go:embed ` directive is configured so that the web server is fully embedded at compile time.

#### 4. Easy to Secure

Only one way in or out, uniform endpoints greatly reduce the surface area you need to think about.

#### 5. No Surprises

No grand abstractions or under-the-hood dynamics. If poetry is not what you are striving for,
a nice bootstrap is a pretty good option!

## QnA

#### Is this a framework or a static site generator?

Neither, it's just a scaffold. Once the code is generated it's yours, adapt it to your needs.
Routing is the only untouchable code, mostly because it is tedious and not very interesting.

#### Why use a special library to manage handlers?

This is what a web framework does that I don't want to live without.
Coordination of requests, handlers and templates needs to be declarative (otherwise 🤯). [Treetop](https://github.com/rur/treetop)
does this with no transitive dependencies.

#### Where does the frontend go?

TODO: ...

#### Is this just a 'code bloat' generator?

Kinda yes, but there is method... Most _reuse_ is multi-purposing, and it sucks.
Reusability (DRY) should emerge from refactoring. Generate some simple bloat-y code, make it work,
refactoring as you go. The patterns will be obvious after a while.
