# _Good_

## A pretty good web app scaffold for Golang

1. Generate obvious, grok-able code
1. Only basic dependencies
1. Binary embedding
1. Easy to secure
1. No surprises

`good` is a code-gen tool for embedding a web GUI to your Golang project. The bootstrap
is geared towards an admin portal or server console, but it's a general purpose setup. Suitable
when a framework would be over-doing it.

### CLI Overview

These commands support a rapid development loop: generate code and refactor for your purposes.

#### Good Scaffold

    $ good scaffold portal

This will create a new app at `[pkg]/portal`, reading the module namespace from the `./go.mod` file.

#### Good Page

    $ good page ./portal mypage

Add a new page 'mypage' to the portal (contains its own route config).

#### Good Routes

    $ good routes ./portal/mypage/routemap.toml

(Re)generate the routing code for the config file. This will overwrite `./portal/mypage/routemap.go`
and output any handlers functions that are missing.

## Intro

#### 1. Generate obvious, grok-able code

We output mostly vanilla Golang code and HTML templates. Code-gen helps us to embrace a little
redundancy so our code can be more static and easier to modify. Go tooling and builtin HTML templates
make it fun to generate and refactor new endpoints as needed.

#### 2. Only Basic Dependencies

With few mandatory dependencies<sup>[1]</sup> it is straightforward to
integrate the libraries that have value for a given project.

#### 3. Binary Embedded

The `//go:embed ` directive is configured so that the web server is fully embedded at compile time.

#### 4. Easy to Secure

One way in and out, uniform endpoints greatly reduce the surface area you have to think about.

#### 5. No Surprises

No grand abstractions or under-the-hood dynamics. If prose is not what you are striving for,
a thoughtful bootstrap is a pretty good option!
